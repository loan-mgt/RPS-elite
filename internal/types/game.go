package types

type Game struct {
	Players   map[string]Player `json:"players"`
	Round     int               `json:"round"`
	AllowMove bool
}
