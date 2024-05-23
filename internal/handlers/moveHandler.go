package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"rcp/elite/internal/senders"
	"rcp/elite/internal/services"
	templatedata "rcp/elite/internal/types/template-data"
	"rcp/elite/internal/utils"
	"time"

	"github.com/gorilla/websocket"
)

type MoveRequest struct {
	Move string
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

	playerMoveData := templatedata.Move{
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
				err = senders.SendMove(opponent.Conn, "opponent", request.Move)
				if err != nil {
					log.Println("Failed to send opponent move:", err)
				}

				err = senders.SendMove(player.Conn, "opponent", opponent.Move)
				if err != nil {
					log.Println("Failed to send opponent move:", err)
				}

				message := "next round in 3s"
				if err := senders.SendMessage(conn, message); err != nil {
					log.Println("Error sending message to player:", err)
				}

				opponent, err := services.GetOpponent(player.Name)
				if err != nil {
					log.Println("Failed getting opponent:", err)
				} else {
					if err := senders.SendMessage(opponent.Conn, message); err != nil {
						log.Println("Error sending message to opponent:", err)
					}
				}

				go func() {
					time.Sleep(3 * time.Second)

					services.SetPlayerMove(player.Name, "")

					services.SetPlayerMove(opponent.Name, "")

					err = senders.ResetMove(opponent.Conn, "player")
					if err != nil {
						log.Println("Failed to send player move:", err)
					}

					err = senders.ResetMove(opponent.Conn, "opponent")
					if err != nil {
						log.Println("Failed to send opponent move:", err)
					}

					err = senders.ResetMove(player.Conn, "player")
					if err != nil {
						log.Println("Failed to send player move:", err)
					}

					err = senders.ResetMove(player.Conn, "opponent")
					if err != nil {
						log.Println("Failed to send opponent move:", err)
					}

					err = senders.SetScore(player.Conn, "opponent", 2)
					if err != nil {
						log.Println("Failed to send opponent score:", err)
					}

					err = senders.SetScore(opponent.Conn, "opponent", 2)
					if err != nil {
						log.Println("Failed to send opponent score:", err)
					}

				}()

			}

		} else {
			messagePlayer := "Waiting for your opponent to make a move"
			messageOpponent := fmt.Sprintf("%s is waiting on you, please select a move", player.Name)

			if err := senders.SendMessage(conn, messagePlayer); err != nil {
				log.Println("Error sending message to player:", err)
			}

			opponent, err := services.GetOpponent(player.Name)
			if err != nil {
				log.Println("Failed getting opponent:", err)
			} else {
				if err := senders.SendMessage(opponent.Conn, messageOpponent); err != nil {
					log.Println("Error sending message to opponent:", err)
				}
			}

		}

	}

	return conn.WriteMessage(websocket.TextMessage, tplBuffer.Bytes())
}
