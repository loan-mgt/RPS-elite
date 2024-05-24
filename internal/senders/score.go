package senders

import (
	"bytes"
	"log"
	templatedata "rcp/elite/internal/types/template-data"
	"rcp/elite/internal/utils"

	"github.com/gorilla/websocket"
)

func SetScore(conn *websocket.Conn, of string, score int) error {
	moveData := templatedata.Score{
		TargetId:  of,
		Score:     score,
		ScoreLoop: make([]int, score),
	}

	var tplBufferOpponent bytes.Buffer
	err := utils.Templates.ExecuteTemplate(&tplBufferOpponent, "score", moveData)
	if err != nil {
		log.Printf("Error executing template score of %s: %v\n", of, err)
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, tplBufferOpponent.Bytes())

}
