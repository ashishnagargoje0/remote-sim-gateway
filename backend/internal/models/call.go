package models

import "time"

type Call struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	PhoneNumber string    `json:"phone_number" gorm:"not null"`
	Duration    int       `json:"duration"` // in seconds
	Status      string    `json:"status" gorm:"default:'pending'"` // pending, connected, failed, ended
	ErrorMsg    string    `json:"error_msg,omitempty"`
	DeviceID    uint      `json:"device_id"`
	Device      Device    `json:"device" gorm:"foreignKey:DeviceID"`
	UserID      uint      `json:"user_id"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
	StartedAt   time.Time `json:"started_at"`
	EndedAt     time.Time `json:"ended_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MakeCallRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	DeviceID    uint   `json:"device_id"`
}

type CallResponse struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Status      string `json:"status"`
}

type CallStatusUpdate struct {
	CallID    uint      `json:"call_id"`
	Status    string    `json:"status"`
	Duration  int       `json:"duration"`
	ErrorMsg  string    `json:"error_msg,omitempty"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
}