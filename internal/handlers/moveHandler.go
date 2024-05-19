package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"rcp/elite/internal/services"
	"rcp/elite/internal/utils"
	"time"

	"github.com/gorilla/websocket"
)

type MoveRequest struct {
	Move string
}

type MoveData struct {
	TargetId string
	Move     *string
}

type Messenger struct {
	Message string
}

func HandleMove(message []byte, conn *websocket.Conn) error {
	var request MoveRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		log.Println("Error parsing message:", err)
		return err
	}

	if request.Move != "rock" && request.Move != "paper" && request.Move != "scissor" {
		return errors.New("invalid move")
	}

	player, err := services.GetPlayerFromConn(conn)
	if err != nil {
		log.Println("Error unable to find player:", err)
		return err
	}

	services.SetPlayerMove(player.Name, request.Move)

	playerMoveData := MoveData{
		TargetId: "player-selected-move",
		Move:     &request.Move,
	}

	var tplBuffer bytes.Buffer
	err = utils.Templates.ExecuteTemplate(&tplBuffer, "move", playerMoveData)
	if err != nil {
		log.Println("Error executing template:", err)
		return err
	}

	opponent, err := services.GetOpponent(player.Name)
	if err != nil {
		log.Println("Unalbe to get oppoent infos:", err)
	} else {

		if opponent.Move != "" {

			opponent, err := services.GetOpponent(player.Name)
			if err != nil {
				log.Println("Failed getting opponent:", err)
			} else {
				err = sendMove(opponent.Conn, "opponent", request.Move)
				if err != nil {
					log.Println("Failed to send opponent move:", err)
				}

				err = sendMove(player.Conn, "opponent", opponent.Move)
				if err != nil {
					log.Println("Failed to send opponent move:", err)
				}

				message := "next round in 3s"
				if err := sendMessage(conn, message); err != nil {
					log.Println("Error sending message to player:", err)
				}

				opponent, err := services.GetOpponent(player.Name)
				if err != nil {
					log.Println("Failed getting opponent:", err)
				} else {
					if err := sendMessage(opponent.Conn, message); err != nil {
						log.Println("Error sending message to opponent:", err)
					}
				}

				go func() {
					time.Sleep(3 * time.Second)

					services.SetPlayerMove(player.Name, "")

					services.SetPlayerMove(opponent.Name, "")

					err = resetMove(opponent.Conn, "player")
					if err != nil {
						log.Println("Failed to send player move:", err)
					}

					err = resetMove(opponent.Conn, "opponent")
					if err != nil {
						log.Println("Failed to send opponent move:", err)
					}

					err = resetMove(player.Conn, "player")
					if err != nil {
						log.Println("Failed to send player move:", err)
					}

					err = resetMove(player.Conn, "opponent")
					if err != nil {
						log.Println("Failed to send opponent move:", err)
					}

				}()

			}

		} else {
			messagePlayer := "Waiting for your opponent to make a move"
			messageOpponent := fmt.Sprintf("%s is waiting on you, please select a move", player.Name)

			if err := sendMessage(conn, messagePlayer); err != nil {
				log.Println("Error sending message to player:", err)
			}

			opponent, err := services.GetOpponent(player.Name)
			if err != nil {
				log.Println("Failed getting opponent:", err)
			} else {
				if err := sendMessage(opponent.Conn, messageOpponent); err != nil {
					log.Println("Error sending message to opponent:", err)
				}
			}

		}

	}

	return conn.WriteMessage(websocket.TextMessage, tplBuffer.Bytes())
}

func sendMessage(conn *websocket.Conn, message string) error {

	messenger := Messenger{
		Message: message,
	}

	var tplBuffer bytes.Buffer
	err := utils.Templates.ExecuteTemplate(&tplBuffer, "messenger", messenger)
	if err != nil {
		log.Println("Error executing template for opponent:", err)
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, tplBuffer.Bytes())

}

func sendMove(conn *websocket.Conn, to string, move string) error {
	moveData := MoveData{
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

func resetMove(conn *websocket.Conn, to string) error {
	moveData := MoveData{
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
