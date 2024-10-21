package websockets

import (
	"encoding/json"
	"log"
	"messaging-service/services"
	"messaging-service/utils"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleConnections handles the WebSocket connection and processes messages
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	// Verify token
	_, ok := utils.ValidateToken(token)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract username from the token
	username := utils.GetUsernameFromToken(token)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("User %s connected to WebSocket", username)

	// Add the user to the WebSocket manager
	WSM.AddClient(username, conn)

	for {
		// Read message from WebSocket connection
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			WSM.RemoveClient(username)
			break
		}

		log.Printf("Received message from user %s: %s", username, string(msg))

		// Parse the message and determine the recipient
		var messageData struct {
			ReceiverUsername string `json:"receiver"`
			Content          string `json:"content"`
		}

		if err := json.Unmarshal(msg, &messageData); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		// Find the receiver's WebSocket connection
		receiverConn, ok := WSM.Clients[messageData.ReceiverUsername]
		if !ok {
			log.Printf("Receiver %s not connected", messageData.ReceiverUsername)
			continue
		}

		// Store the message in the database
		sender, err := services.GetUser(username)
		if err != nil {
			log.Printf("Sender user not found: %v", err)
			continue
		}

		receiver, err := services.GetUser(messageData.ReceiverUsername)
		if err != nil {
			log.Printf("Receiver user not found: %v", err)
			continue
		}

		err = services.CreateMessage(sender.ID, receiver.ID, messageData.Content)
		if err != nil {
			log.Printf("Error storing message in database: %v", err)
			continue
		}
		// Parse the message and determine the recipient
		var Mesg struct {
			Sender  string `json:"sender"`
			Content string `json:"content"`
		}
		Mesg.Sender = sender.Username
		Mesg.Content = messageData.Content

		rcvMsg, err := json.Marshal(Mesg)
		if err != nil {
			log.Printf("Error sending message11 to %s: %v", messageData.ReceiverUsername, err)
			continue
		}
		// Send the message to the recipient WebSocket
		if err := receiverConn.WriteMessage(websocket.TextMessage, rcvMsg); err != nil {
			log.Printf("Error sending message to %s: %v", messageData.ReceiverUsername, err)
			continue
		}
	}
}
