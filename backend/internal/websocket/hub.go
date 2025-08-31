package websocket

import (
	"log"
	"sync"
	"time"
)

type Hub struct {
	// Registered clients
	Clients map[*Client]bool

	// Map device ID to client for direct messaging
	DeviceMap map[string]*Client

	// Inbound messages from clients
	Broadcast chan Message

	// Register requests from clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client

	// Mutex for thread-safe operations
	mutex sync.RWMutex

	// Device status tracking
	DeviceStatus map[string]*DeviceInfo
}

type DeviceInfo struct {
	DeviceID       string    `json:"device_id"`
	IsOnline       bool      `json:"is_online"`
	LastSeen       time.Time `json:"last_seen"`
	BatteryLevel   int       `json:"battery_level"`
	SignalStrength int       `json:"signal_strength"`
	PhoneNumber    string    `json:"phone_number"`
}

type Message struct {
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
	DeviceID  string                 `json:"device_id,omitempty"`
}

func NewHub() *Hub {
	return &Hub{
		Clients:      make(map[*Client]bool),
		DeviceMap:    make(map[string]*Client),
		Broadcast:    make(chan Message, 256),
		Register:     make(chan *Client),
		Unregister:   make(chan *Client),
		DeviceStatus: make(map[string]*DeviceInfo),
	}
}

func (h *Hub) Run() {
	// Start cleanup routine
	go h.startCleanupRoutine()

	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)

		case client := <-h.Unregister:
			h.unregisterClient(client)

		case message := <-h.Broadcast:
			h.broadcastMessage(message)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.Clients[client] = true
	if client.DeviceID != "" {
		h.DeviceMap[client.DeviceID] = client
		
		// Update device status
		h.DeviceStatus[client.DeviceID] = &DeviceInfo{
			DeviceID: client.DeviceID,
			IsOnline: true,
			LastSeen: time.Now(),
		}
	}

	log.Printf("Client registered: %s (Total clients: %d)", client.DeviceID, len(h.Clients))

	// Send welcome message
	welcomeMsg := Message{
		Type:      "welcome",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message":   "Connected to Remote SIM Gateway",
			"device_id": client.DeviceID,
			"server_time": time.Now().Unix(),
		},
	}

	select {
	case client.Send <- welcomeMsg:
	default:
		log.Printf("Failed to send welcome message to %s", client.DeviceID)
		close(client.Send)
		delete(h.Clients, client)
		if client.DeviceID != "" {
			delete(h.DeviceMap, client.DeviceID)
			delete(h.DeviceStatus, client.DeviceID)
		}
	}
}

func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.Clients[client]; ok {
		delete(h.Clients, client)
		
		if client.DeviceID != "" {
			delete(h.DeviceMap, client.DeviceID)
			
			// Update device status to offline
			if deviceInfo, exists := h.DeviceStatus[client.DeviceID]; exists {
				deviceInfo.IsOnline = false
				deviceInfo.LastSeen = time.Now()
			}
		}
		
		close(client.Send)
		log.Printf("Client unregistered: %s (Total clients: %d)", client.DeviceID, len(h.Clients))
	}
}

func (h *Hub) broadcastMessage(message Message) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	message.Timestamp = time.Now()
	
	for client := range h.Clients {
		select {
		case client.Send <- message:
		default:
			log.Printf("Failed to send broadcast message to %s", client.DeviceID)
			close(client.Send)
			delete(h.Clients, client)
			if client.DeviceID != "" {
				delete(h.DeviceMap, client.DeviceID)
			}
		}
	}
	
	log.Printf("Broadcast message sent to %d clients: %s", len(h.Clients), message.Type)
}

func (h *Hub) SendToDevice(deviceID string, message Message) bool {
	h.mutex.RLock()
	client, ok := h.DeviceMap[deviceID]
	h.mutex.RUnlock()

	if !ok {
		log.Printf("Device not found or offline: %s", deviceID)
		return false
	}

	message.Timestamp = time.Now()
	message.DeviceID = deviceID

	select {
	case client.Send <- message:
		log.Printf("Message sent to device %s: %s", deviceID, message.Type)
		return true
	default:
		log.Printf("Failed to send message to device %s: channel full", deviceID)
		return false
	}
}

func (h *Hub) GetConnectedDevices() []string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	var devices []string
	for deviceID := range h.DeviceMap {
		if deviceID != "" {
			devices = append(devices, deviceID)
		}
	}
	return devices
}

func (h *Hub) GetDeviceStatus(deviceID string) (*DeviceInfo, bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	deviceInfo, exists := h.DeviceStatus[deviceID]
	return deviceInfo, exists
}

func (h *Hub) UpdateDeviceStatus(deviceID string, update map[string]interface{}) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	deviceInfo, exists := h.DeviceStatus[deviceID]
	if !exists {
		deviceInfo = &DeviceInfo{
			DeviceID: deviceID,
			IsOnline: false,
		}
		h.DeviceStatus[deviceID] = deviceInfo
	}

	// Update fields from the update map
	if batteryLevel, ok := update["battery_level"].(float64); ok {
		deviceInfo.BatteryLevel = int(batteryLevel)
	}
	
	if signalStrength, ok := update["signal_strength"].(float64); ok {
		deviceInfo.SignalStrength = int(signalStrength)
	}
	
	if phoneNumber, ok := update["phone_number"].(string); ok {
		deviceInfo.PhoneNumber = phoneNumber
	}

	deviceInfo.LastSeen = time.Now()
	deviceInfo.IsOnline = true

	log.Printf("Device status updated: %s", deviceID)
}

func (h *Hub) GetHubStats() map[string]interface{} {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	onlineDevices := 0
	for _, deviceInfo := range h.DeviceStatus {
		if deviceInfo.IsOnline {
			onlineDevices++
		}
	}

	return map[string]interface{}{
		"total_clients":    len(h.Clients),
		"online_devices":   onlineDevices,
		"total_devices":    len(h.DeviceStatus),
		"uptime":          time.Since(time.Now()).String(), // This would be set when hub starts
	}
}

func (h *Hub) startCleanupRoutine() {
	ticker := time.NewTicker(5 * time.Minute) // Clean up every 5 minutes
	defer ticker.Stop()

	for range ticker.C {
		h.cleanupStaleDevices()
	}
}

func (h *Hub) cleanupStaleDevices() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	staleThreshold := time.Now().Add(-10 * time.Minute) // 10 minutes

	for deviceID, deviceInfo := range h.DeviceStatus {
		if deviceInfo.LastSeen.Before(staleThreshold) && deviceInfo.IsOnline {
			deviceInfo.IsOnline = false
			log.Printf("Marked device as offline due to inactivity: %s", deviceID)
		}
	}
}

// SendHeartbeat sends a heartbeat to all connected devices
func (h *Hub) SendHeartbeat() {
	heartbeatMsg := Message{
		Type: "heartbeat",
		Data: map[string]interface{}{
			"timestamp": time.Now().Unix(),
		},
	}

	h.Broadcast <- heartbeatMsg
}

// StartHeartbeat starts a routine to send periodic heartbeats
func (h *Hub) StartHeartbeat() {
	ticker := time.NewTicker(30 * time.Second) // Send heartbeat every 30 seconds
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			h.SendHeartbeat()
		}
	}()
}