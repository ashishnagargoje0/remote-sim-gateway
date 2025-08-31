package models

import "time"

type Device struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	DeviceID     string    `json:"device_id" gorm:"unique;not null"`
	Name         string    `json:"name"`
	PhoneNumber  string    `json:"phone_number"`
	IsOnline     bool      `json:"is_online" gorm:"default:false"`
	LastSeenAt   time.Time `json:"last_seen_at"`
	UserID       uint      `json:"user_id"`
	User         User      `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type DeviceRegisterRequest struct {
	DeviceID    string `json:"device_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number"`
}

type DeviceUpdateRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	IsOnline    bool   `json:"is_online"`
}

type DeviceStatusUpdate struct {
	DeviceID     string    `json:"device_id"`
	IsOnline     bool      `json:"is_online"`
	BatteryLevel int       `json:"battery_level"`
	SignalStrength int     `json:"signal_strength"`
	LastSeenAt   time.Time `json:"last_seen_at"`
}