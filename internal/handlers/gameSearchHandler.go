package handlers

import (
	"encoding/json"
	"log"

	"rcp/elite/internal/senders"
	"rcp/elite/internal/services"
	"rcp/elite/internal/types"

	"github.com/gorilla/websocket"
)

type GameSearchRequest struct {
	Username string `json:"username"`
}

func HandleGameSearch(message []byte, conn *websocket.Conn) (string, error) {
	var request GameSearchRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		log.Println("Error parsing message:", err)
		return "", err
	}

	player := &types.Player{
		Name:  request.Username,
		Move:  "",
		Flag:  "FR",
		Score: 0,
		Conn:  conn,
	}

	services.JoinPlayerPoll(*player)

	return player.Name, senders.SendWaitScreen(player.Conn)
}
