package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"rcp/elite/internal/utils"
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

	var tplBuffer bytes.Buffer
	err = utils.Templates.ExecuteTemplate(&tplBuffer, "move", request)
	if err != nil {
		log.Println("Error executing template:", err)
		return nil, err
	}

	return tplBuffer.Bytes(), nil
}
