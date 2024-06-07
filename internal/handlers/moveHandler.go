package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"rcp/elite/internal/senders"
	"rcp/elite/internal/services"

	"github.com/gorilla/websocket"
)

type MoveRequest struct {
	Move string
}

func HandleMove(message []byte, playerName string, conn *websocket.Conn) error {
	var request MoveRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		log.Println("Error parsing message:", err)
		return err
	}

	if request.Move != "rock" && request.Move != "paper" && request.Move != "scissor" {
		return errors.New("invalid move")
	}

	err = services.SetPlayerMove(playerName, request.Move)

	if err != nil {
		return nil
	}
	senders.SendMove(conn, "player", request.Move)

	return nil

}
