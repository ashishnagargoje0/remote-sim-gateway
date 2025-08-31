package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"remote-sim-gateway/internal/models"
	"remote-sim-gateway/internal/websocket"
)

type DeviceHandler struct {
	db  *gorm.DB
	hub *websocket.Hub
}

func NewDeviceHandler(db *gorm.DB, hub *websocket.Hub) *DeviceHandler {
	return &DeviceHandler{
		db:  db,
		hub: hub,
	}
}

func (h *DeviceHandler) GetDevices(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var devices []models.Device
	if err := h.db.Where("user_id = ?", userID).Find(&devices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch devices"})
		return
	}

	// Update online status based on WebSocket connections
	connectedDevices := h.hub.GetConnectedDevices()
	for i := range devices {
		devices[i].IsOnline = false
		for _, connectedID := range connectedDevices {
			if devices[i].DeviceID == connectedID {
				devices[i].IsOnline = true
				devices[i].LastSeenAt = time.Now()
				break
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"devices": devices})
}

func (h *DeviceHandler) RegisterDevice(c *gin.Context) {
	var req models.DeviceRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	// Check if device already exists
	var existingDevice models.Device
	if err := h.db.Where("device_id = ?", req.DeviceID).First(&existingDevice).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Device already registered"})
		return
	}

	device := models.Device{
		DeviceID:    req.DeviceID,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		UserID:      userID.(uint),
		IsOnline:    false,
		LastSeenAt:  time.Now(),
	}

	if err := h.db.Create(&device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register device"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Device registered successfully",
		"device":  device,
	})
}

func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
	deviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	var req models.DeviceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	var device models.Device
	if err := h.db.Where("id = ? AND user_id = ?", deviceID, userID).First(&device).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	// Update device fields
	if req.Name != "" {
		device.Name = req.Name
	}
	if req.PhoneNumber != "" {
		device.PhoneNumber = req.PhoneNumber
	}

	if err := h.db.Save(&device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update device"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Device updated successfully",
		"device":  device,
	})
}

func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	deviceID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	userID, _ := c.Get("user_id")

	var device models.Device
	if err := h.db.Where("id = ? AND user_id = ?", deviceID, userID).First(&device).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	if err := h.db.Delete(&device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete device"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device deleted successfully"})
}