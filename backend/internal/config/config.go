package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

type JWTConfig struct {
	Secret         string
	ExpirationTime int
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

type ServerConfig struct {
	Port            string
	ReadTimeout     int
	WriteTimeout    int
	ShutdownTimeout int
}

func New() *Config {
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	jwtExpiration, _ := strconv.Atoi(getEnv("JWT_EXPIRATION", "24"))

	origins := strings.Split(getEnv("CORS_ORIGINS", "http://localhost:3000"), ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			Name:     getEnv("DB_NAME", "simgateway"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:         getEnv("JWT_SECRET", "your-super-secret-key"),
			ExpirationTime: jwtExpiration,
		},
		CORS: CORSConfig{
			AllowedOrigins: origins,
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type", "Authorization"},
		},
		Server: ServerConfig{
			Port:            getEnv("PORT", "8080"),
			ReadTimeout:     30,
			WriteTimeout:    30,
			ShutdownTimeout: 10,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}