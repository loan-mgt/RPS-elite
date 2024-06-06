package senders

import (
	"bytes"
	"log"
	"rcp/elite/internal/utils"

	"github.com/gorilla/websocket"
)

func SendWaitScreen(conn *websocket.Conn) error {

	var tplBuffer bytes.Buffer
	err := utils.Templates.ExecuteTemplate(&tplBuffer, "wait", nil)
	if err != nil {
		log.Println("Error executing template for wait:", err)
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, tplBuffer.Bytes())

}
