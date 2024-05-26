package senders

import (
	"bytes"
	"log"
	templatedata "rcp/elite/internal/types/template-data"
	"rcp/elite/internal/utils"

	"github.com/gorilla/websocket"
)

func SendMessage(conn *websocket.Conn, message string, timer *string) error {

	messenger := templatedata.Messenger{
		Message: message,
		Timer:   timer,
	}

	var tplBuffer bytes.Buffer
	err := utils.Templates.ExecuteTemplate(&tplBuffer, "messenger", messenger)
	if err != nil {
		log.Println("Error executing template for opponent:", err)
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, tplBuffer.Bytes())

}
