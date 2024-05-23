package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"

	"rcp/elite/internal/services"
	"rcp/elite/internal/types"
	templatedata "rcp/elite/internal/types/template-data"
	"rcp/elite/internal/utils"

	"github.com/gorilla/websocket"
)

type GameSearchRequest struct {
	Username string `json:"username"`
}

func HandleGameSearch(message []byte, conn *websocket.Conn) error {
	var request GameSearchRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		log.Println("Error parsing message:", err)
		return err
	}

	player := &types.Player{
		Name:  request.Username,
		Move:  "",
		Flag:  "FR",
		Score: 0,
		Conn:  conn,
	}

	if !services.IsPlayerInGame(request.Username) && services.IsGameFull() {
		return errors.New("no game available")
	} else if !services.IsPlayerInGame(request.Username) {
		err := services.AddPlayer(player)
		if err != nil {
			log.Println("Error adding player:", err)
			return err
		}
	}

	opponent, err := services.GetOpponent(request.Username)
	if err != nil {
		log.Println("Error getting opponent:", err)
	} else {
		log.Println("Opponent:", opponent)

		opponentInfo := templatedata.OpponentInfo{
			Opponent: player,
			Score: templatedata.Score{
				TargetId: "opponent",
				Score:    0,
			},
		}

		var tplBuffer bytes.Buffer
		err = utils.Templates.ExecuteTemplate(&tplBuffer, "opponent-info", opponentInfo)
		if err != nil {
			log.Println("Error executing template:", err)
			return err
		}

		opponent.Conn.WriteMessage(websocket.TextMessage, tplBuffer.Bytes())

		var tplBuffer2 bytes.Buffer
		err = utils.Templates.ExecuteTemplate(&tplBuffer2, "opponent-panel", opponentInfo)
		if err != nil {
			log.Println("Error executing template:", err)
			return err
		}

		opponent.Conn.WriteMessage(websocket.TextMessage, tplBuffer2.Bytes())
	}

	players := templatedata.Home{
		Player: player,
		OpponentInfo: &templatedata.OpponentInfo{
			Opponent: opponent,
			Score: templatedata.Score{
				TargetId: "opponent",
				Score:    0,
			},
		},
		Messenger: templatedata.Messenger{
			Message: "Welcome",
		},
	}

	// Parse the template
	var tplBuffer bytes.Buffer
	err = utils.Templates.ExecuteTemplate(&tplBuffer, "gameHome", players)
	if err != nil {
		log.Println("Error executing template:", err)
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, tplBuffer.Bytes())
}
