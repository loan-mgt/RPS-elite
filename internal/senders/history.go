package senders

import (
	"bytes"
	"log"
	templatedata "rcp/elite/internal/types/template-data"
	"rcp/elite/internal/utils"

	"github.com/gorilla/websocket"
)

func AppendHistory(conn *websocket.Conn, playerMove, opponentMove, winner string) error {
	historyData := templatedata.History{
		PlayerMove:   playerMove,
		OpponentMove: opponentMove,
		Winner:       winner,
	}

	var tplBufferOpponent bytes.Buffer
	err := utils.Templates.ExecuteTemplate(&tplBufferOpponent, "history-row", historyData)
	if err != nil {
		log.Printf("Error executing template history : %v\n", err)
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, tplBufferOpponent.Bytes())

}
