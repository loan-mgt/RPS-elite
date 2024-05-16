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

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("Received message: %s type: %s\n", message, string(rune(messageType)))

		var response []byte

		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("Error parsing message:", err)
			continue
		}

		switch msg.Type {
		case "game-search":
			response, err = handlers.HandleGameSearch(message)
		case "move":
			response, err = handlers.HandleMove(message)
		default:
			log.Println("Unknown message type:", msg.Type)
			continue
		}

		if err != nil {
			log.Println("Error writing response to client:", err)
			continue
		}

		err = conn.WriteMessage(websocket.TextMessage, response)
		if err != nil {
			log.Println("Error writing response to client:", err)
			continue
		}
	}
}
