package senders

import (
	"bytes"
	"log"
	"rcp/elite/internal/types"
	templatedata "rcp/elite/internal/types/template-data"
	"rcp/elite/internal/utils"

	"github.com/gorilla/websocket"
)

func SetGameHome(conn *websocket.Conn, player, opponent *types.Player, messsage, timer string) error {
	gameHomeData := templatedata.GameHome{
		PlayerInfo: &templatedata.PlayerInfo{
			Player: player,
			Score: templatedata.Score{
				TargetId:  "player",
				Score:     player.Score,
				ScoreLoop: make([]int, player.Score),
			},
		},
		OpponentInfo: &templatedata.PlayerInfo{
			Player: opponent,
			Score: templatedata.Score{
				TargetId:  "opponent",
				Score:     opponent.Score,
				ScoreLoop: make([]int, opponent.Score),
			},
		},
		Messenger: templatedata.Messenger{
			Message: messsage,
			Timer:   timer,
		},
	}

	var tplBufferOpponent bytes.Buffer
	err := utils.Templates.ExecuteTemplate(&tplBufferOpponent, "gameHome", gameHomeData)
	if err != nil {
		log.Printf("Error executing template gameHome of %s: %v\n", player.Name, err)
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, tplBufferOpponent.Bytes())

}
