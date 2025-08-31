package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"remote-sim-gateway/internal/config"
	"remote-sim-gateway/internal/database"
	"remote-sim-gateway/internal/handlers"
	"remote-sim-gateway/internal/middleware"
	"remote-sim-gateway/internal/websocket"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize configuration
	cfg := config.New()

	// Initialize database
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize Gin router
	router := gin.Default()

	// Apply middleware
	router.Use(middleware.CORS(cfg.CORS))
	router.Use(middleware.Logger())
	router.Use(middleware.RateLimit())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg.JWT)
	smsHandler := handlers.NewSMSHandler(db, hub)
	callHandler := handlers.NewCallHandler(db, hub)
	deviceHandler := handlers.NewDeviceHandler(db, hub)
	dashboardHandler := handlers.NewDashboardHandler(db)

	// Public routes
	public := router.Group("/")
	{
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/register", authHandler.Register)
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.AuthRequired(cfg.JWT.Secret))
	{
		// SMS routes
		api.POST("/send-sms", smsHandler.SendSMS)
		api.POST("/send-bulk-sms", smsHandler.SendBulkSMS)
		api.GET("/sms-history", smsHandler.GetHistory)

		// Call routes
		api.POST("/make-call", callHandler.MakeCall)
		api.GET("/call-history", callHandler.GetHistory)

		// Device routes
		api.GET("/devices", deviceHandler.GetDevices)
		api.POST("/devices", deviceHandler.RegisterDevice)
		api.PUT("/devices/:id", deviceHandler.UpdateDevice)
		api.DELETE("/devices/:id", deviceHandler.DeleteDevice)

		// Dashboard routes
		api.GET("/dashboard/stats", dashboardHandler.GetStats)
		api.GET("/dashboard/activity", dashboardHandler.GetRecentActivity)
	}

	// WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		websocket.HandleWebSocket(hub, c.Writer, c.Request)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}