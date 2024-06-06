package handlers

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gorilla/websocket"
)

type MoveRequest struct {
	Move string
}

func HandleMove(message []byte, conn *websocket.Conn) error {
	var request MoveRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		log.Println("Error parsing message:", err)
		return err
	}

	if request.Move != "rock" && request.Move != "paper" && request.Move != "scissor" {
		return errors.New("invalid move")
	}

	return nil

}
