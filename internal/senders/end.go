package senders

import (
	"bytes"
	"log"
	templatedata "rcp/elite/internal/types/template-data"
	"rcp/elite/internal/utils"

	"github.com/gorilla/websocket"
)

func SendEndScreen(conn *websocket.Conn, message string) error {

	end := templatedata.End{
		Message: message,
	}

	var tplBuffer bytes.Buffer
	err := utils.Templates.ExecuteTemplate(&tplBuffer, "end", end)
	if err != nil {
		log.Println("Error executing template for end:", err)
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, tplBuffer.Bytes())

}
