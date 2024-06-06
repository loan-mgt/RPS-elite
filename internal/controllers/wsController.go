package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"rcp/elite/internal/handlers"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Request struct {
	Type string `json:"type"`
}

type Message struct {
	Type    string `json:"type"`
	Success bool   `json:"success,omitempty"`
}

func MainController(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	handleWebSocket(conn)
}

func handleWebSocket(conn *websocket.Conn) {
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			handleReadError(err, conn)
			break
		}

		log.Printf("Received message: %s type: %s\n", message, string(rune(messageType)))

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("Error parsing message:", err)
			continue
		}

		err = handleMessage(msg, message, conn)
		if err != nil {
			handleWriteError(err, conn)
			continue
		}
	}
}

func handleReadError(err error, conn *websocket.Conn) {
	log.Printf("Error reading message: %v", err)

}

func handleWriteError(err error, conn *websocket.Conn) {
	log.Printf("Error writing response to client: %v", err)
}

func handleMessage(msg Message, originalMessage []byte, conn *websocket.Conn) error {
	switch msg.Type {
	case "game-search":
		return handlers.HandleGameSearch(originalMessage, conn)
	case "move":
		return handlers.HandleMove(originalMessage, conn)
	default:
		log.Println("Unknown message type:", msg.Type)
		return nil
	}
}
