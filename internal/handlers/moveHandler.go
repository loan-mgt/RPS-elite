package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type MoveRequest struct {
	Move string `json:"move"`
}

func HandleMove(message []byte) ([]byte, error) {
	var request MoveRequest
	err := json.Unmarshal(message, &request)
	if err != nil {
		log.Println("Error parsing message:", err)
		return nil, err
	}

	if request.Move != "rock" && request.Move != "paper" && request.Move != "scissor" {
		return nil, errors.New("invalid move")
	}

	content := fmt.Sprintf(`
		<div hx-swap-oob="innerHTML:#player-selected-move" >
			<img class="w-full h-full" src="/assets/images/%s.svg" alt="Rock" />
		</div>
	`, request.Move)

	return []byte(content), nil
}
