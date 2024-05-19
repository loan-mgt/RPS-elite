package types

type Game struct {
	Player1 *Player `json:"player1"`
	Player2 *Player `json:"player2"`
	Round   int     `json:"round"`
}
