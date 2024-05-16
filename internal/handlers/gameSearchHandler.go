package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"

	"github.com/gorilla/websocket"
)

type Player struct {
	ID    string          `json:"id"`
	Name  string          `json:"name"`
	Move  string          `json:"move"`
	Flag  string          `json:"flag"`
	Score int             `json:"score"`
	Conn  *websocket.Conn `json:"conn"`
}

func (p Player) String() string {
	return fmt.Sprintf("Player{ID: %s, Name: %s, Move: %s, Flag: %s, Score: %d}", p.ID, p.Name, p.Move, p.Flag, p.Score)
}

type PlayersData struct {
	Player   Player
	Opponent Player
}

type GameSearchRequest struct {
	Username string `json:"username"`
}

func HandleGameSearch(message []byte) ([]byte, error) {

	var request GameSearchRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		log.Println("Error parsing message:", err)
		return nil, err
	}

	player1 := Player{
		ID:    "1",
		Name:  "Player 1",
		Move:  "rock",
		Flag:  "BE",
		Score: 0,
		Conn:  nil, // Add the appropriate connection
	}

	player2 := Player{
		ID:    "2",
		Name:  request.Username,
		Move:  "paper",
		Flag:  "FR",
		Score: 0,
		Conn:  nil, // Add the appropriate connection
	}

	// Create a struct containing both players
	players := PlayersData{
		Player:   player1,
		Opponent: player2,
	}

	// Parse the template
	tmpl, err := template.ParseFiles("internal/component/gameHome.gohtml")
	if err != nil {
		log.Println("Error parsing template:", err)
		return nil, err
	}

	// Execute the template with the players data
	var tplBuffer bytes.Buffer
	err = tmpl.Execute(&tplBuffer, players)
	if err != nil {
		log.Println("Error executing template:", err)
		return nil, err
	}

	return tplBuffer.Bytes(), nil

}
