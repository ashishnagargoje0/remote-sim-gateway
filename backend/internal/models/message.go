package models

import "time"

type Message struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	PhoneNumber string    `json:"phone_number" gorm:"not null"`
	Content     string    `json:"content" gorm:"not null"`
	Status      string    `json:"status" gorm:"default:'pending'"` // pending, sent, failed
	ErrorMsg    string    `json:"error_msg,omitempty"`
	DeviceID    uint      `json:"device_id"`
	Device      Device    `json:"device" gorm:"foreignKey:DeviceID"`
	UserID      uint      `json:"user_id"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
	SentAt      time.Time `json:"sent_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SendSMSRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Message     string `json:"message" binding:"required"`
	DeviceID    uint   `json:"device_id"`
}

type BulkSMSRequest struct {
	PhoneNumbers []string `json:"phone_numbers" binding:"required"`
	Message      string   `json:"message" binding:"required"`
	DeviceID     uint     `json:"device_id"`
}

type SMSResponse struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Status      string `json:"status"`
	Message     string `json:"message,omitempty"`
}

type SMSStatusUpdate struct {
	MessageID uint   `json:"message_id"`
	Status    string `json:"status"`
	ErrorMsg  string `json:"error_msg,omitempty"`
	SentAt    time.Time `json:"sent_at"`
}