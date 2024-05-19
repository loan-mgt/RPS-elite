package types

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Player struct {
	Name  string          `json:"name"`
	Move  string          `json:"move"`
	Flag  string          `json:"flag"`
	Score int             `json:"score"`
	Conn  *websocket.Conn `json:"conn"`
}

func (p Player) String() string {
	return fmt.Sprintf("Player{Name: %s, Move: %s, Flag: %s, Score: %d}", p.Name, p.Move, p.Flag, p.Score)
}
