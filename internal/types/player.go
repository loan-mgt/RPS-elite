package types

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Player struct {
	Name  string          `json:"name"`
	Move  *string         `json:"move"`
	Flag  string          `json:"flag"`
	Score int             `json:"score"`
	Conn  *websocket.Conn `json:"conn"`
}

func (p Player) String() string {
	move := "nil"
	if p.Move != nil {
		move = *p.Move
	}
	return fmt.Sprintf("Player{Name: %s, Move: %s, Flag: %s, Score: %d}", p.Name, move, p.Flag, p.Score)
}
