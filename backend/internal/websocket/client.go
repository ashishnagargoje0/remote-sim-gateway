package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan Message
	DeviceID string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implement proper origin checking based on configuration
		return true
	},
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	deviceID := r.URL.Query().Get("device_id")
	if deviceID == "" {
		log.Printf("Missing device_id parameter")
		conn.Close()
		return
	}

	client := &Client{
		Hub:      hub,
		Conn:     conn,
		Send:     make(chan Message, 256),
		DeviceID: deviceID,
	}

	client.Hub.Register <- client

	// Start goroutines for handling read/write
	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var message Message
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			log.Printf("JSON unmarshal error: %v", err)
			continue
		}

		// Handle incoming messages from Android device
		go c.handleMessage(message)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(message Message) {
	log.Printf("Received message from device %s: %s", c.DeviceID, message.Type)

	switch message.Type {
	case "sms_status":
		c.handleSMSStatus(message)
	case "call_status":
		c.handleCallStatus(message)
	case "device_status":
		c.handleDeviceStatus(message)
	case "heartbeat":
		c.handleHeartbeat(message)
	default:
		log.Printf("Unknown message type: %s", message.Type)
	}
}

func (c *Client) handleSMSStatus(message Message) {
	// Update SMS status in database
	log.Printf("SMS status update: %+v", message.Data)
	// TODO: Implement database update logic
}

func (c *Client) handleCallStatus(message Message) {
	// Update call status in database
	log.Printf("Call status update: %+v", message.Data)
	// TODO: Implement database update logic
}

func (c *Client) handleDeviceStatus(message Message) {
	// Update device status in database
	log.Printf("Device status update: %+v", message.Data)
	// TODO: Implement database update logic
}

func (c *Client) handleHeartbeat(message Message) {
	// Send heartbeat response
	response := Message{
		Type: "heartbeat_ack",
		Data: map[string]interface{}{
			"timestamp": time.Now().Unix(),
		},
	}
	
	select {
	case c.Send <- response:
	default:
		log.Printf("Failed to send heartbeat ack to device: %s", c.DeviceID)
	}
}