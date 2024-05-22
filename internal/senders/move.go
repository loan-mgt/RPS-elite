package senders

import (
	"bytes"
	"fmt"
	"log"
	templatedata "rcp/elite/internal/types/template-data"
	"rcp/elite/internal/utils"

	"github.com/gorilla/websocket"
)

func SendMove(conn *websocket.Conn, to string, move string) error {
	moveData := templatedata.Move{
		TargetId: fmt.Sprintf("%s-selected-move", to),
		Move:     &move,
	}

	var tplBufferOpponent bytes.Buffer
	err := utils.Templates.ExecuteTemplate(&tplBufferOpponent, "move", moveData)
	if err != nil {
		log.Println("Error executing template for opponent:", err)
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, tplBufferOpponent.Bytes())

}

func ResetMove(conn *websocket.Conn, to string) error {
	moveData := templatedata.Move{
		TargetId: fmt.Sprintf("%s-selected-move", to),
	}

	var tplBufferOpponent bytes.Buffer
	err := utils.Templates.ExecuteTemplate(&tplBufferOpponent, "move", moveData)
	if err != nil {
		log.Println("Error executing template for opponent:", err)
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, tplBufferOpponent.Bytes())

}
