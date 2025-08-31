package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"remote-sim-gateway/internal/models"
	"remote-sim-gateway/internal/websocket"
)

type SMSHandler struct {
	db  *gorm.DB
	hub *websocket.Hub
}

func NewSMSHandler(db *gorm.DB, hub *websocket.Hub) *SMSHandler {
	return &SMSHandler{
		db:  db,
		hub: hub,
	}
}

func (h *SMSHandler) SendSMS(c *gin.Context) {
	var req models.SendSMSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	
	// Get available device if not specified
	if req.DeviceID == 0 {
		var device models.Device
		if err := h.db.Where("user_id = ? AND is_online = ?", userID, true).First(&device).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No online device available"})
			return
		}
		req.DeviceID = device.ID
	}

	// Verify device belongs to user and is online
	var device models.Device
	if err := h.db.Where("id = ? AND user_id = ? AND is_online = ?", req.DeviceID, userID, true).First(&device).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device not found or offline"})
		return
	}

	// Create message record
	message := models.Message{
		PhoneNumber: req.PhoneNumber,
		Content:     req.Message,
		Status:      "pending",
		DeviceID:    req.DeviceID,
		UserID:      userID.(uint),
	}

	if err := h.db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
		return
	}

	// Send to device via WebSocket
	wsMessage := websocket.Message{
		Type: "send_sms",
		Data: map[string]interface{}{
			"id":           message.ID,
			"phone_number": req.PhoneNumber,
			"message":      req.Message,
			"device_id":    device.DeviceID,
		},
	}

	h.hub.SendToDevice(device.DeviceID, wsMessage)

	c.JSON(http.StatusOK, models.SMSResponse{
		ID:          message.ID,
		PhoneNumber: req.PhoneNumber,
		Status:      "pending",
	})
}

func (h *SMSHandler) SendBulkSMS(c *gin.Context) {
	var req models.BulkSMSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	var responses []models.SMSResponse

	// Get available device if not specified
	if req.DeviceID == 0 {
		var device models.Device
		if err := h.db.Where("user_id = ? AND is_online = ?", userID, true).First(&device).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No online device available"})
			return
		}
		req.DeviceID = device.ID
	}

	// Verify device belongs to user and is online
	var device models.Device
	if err := h.db.Where("id = ? AND user_id = ? AND is_online = ?", req.DeviceID, userID, true).First(&device).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device not found or offline"})
		return
	}

	for _, phoneNumber := range req.PhoneNumbers {
		message := models.Message{
			PhoneNumber: phoneNumber,
			Content:     req.Message,
			Status:      "pending",
			DeviceID:    req.DeviceID,
			UserID:      userID.(uint),
		}

		if err := h.db.Create(&message).Error; err != nil {
			responses = append(responses, models.SMSResponse{
				PhoneNumber: phoneNumber,
				Status:      "failed",
				Message:     "Failed to create message",
			})
			continue
		}

		// Send to device via WebSocket
		wsMessage := websocket.Message{
			Type: "send_sms",
			Data: map[string]interface{}{
				"id":           message.ID,
				"phone_number": phoneNumber,
				"message":      req.Message,
				"device_id":    device.DeviceID,
			},
		}

		h.hub.SendToDevice(device.DeviceID, wsMessage)

		responses = append(responses, models.SMSResponse{
			ID:          message.ID,
			PhoneNumber: phoneNumber,
			Status:      "pending",
		})
	}

	c.JSON(http.StatusOK, gin.H{"messages": responses})
}

func (h *SMSHandler) GetHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	phoneNumber := c.Query("phone_number")
	
	offset := (page - 1) * limit
	userID, _ := c.Get("user_id")

	var messages []models.Message
	var total int64

	query := h.db.Where("user_id = ?", userID).Preload("Device")
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	if phoneNumber != "" {
		query = query.Where("phone_number ILIKE ?", "%"+phoneNumber+"%")
	}
	
	if err := query.Model(&models.Message{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count messages"})
		return
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}