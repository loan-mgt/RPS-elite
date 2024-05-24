package services

import (
	"errors"
	"rcp/elite/internal/types"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	gameInstance = &types.Game{}
	mutex        = &sync.Mutex{}
)

func IsGameFull() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return gameInstance.Player1 != nil && gameInstance.Player2 != nil
}

func IsPlayerInGame(playerName string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	return (gameInstance.Player1 != nil && gameInstance.Player1.Name == playerName) ||
		(gameInstance.Player2 != nil && gameInstance.Player2.Name == playerName)
}

func AddPlayer(player *types.Player) error {
	mutex.Lock()
	defer mutex.Unlock()

	if gameInstance.Player1 == nil {
		gameInstance.Player1 = player
		return nil
	} else if gameInstance.Player2 == nil {
		gameInstance.Player2 = player
		return nil
	} else {
		return errors.New("game is already full")
	}
}

func GetPlayerStatus(playerName string) (*types.Player, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if gameInstance.Player1 != nil && gameInstance.Player1.Name == playerName {
		return gameInstance.Player1, nil
	} else if gameInstance.Player2 != nil && gameInstance.Player2.Name == playerName {
		return gameInstance.Player2, nil
	} else {
		return nil, errors.New("player not found")
	}
}

func SetPlayerMove(playerName, move string) error {
	mutex.Lock()
	defer mutex.Unlock()

	if gameInstance.Player1 != nil && gameInstance.Player1.Name == playerName {
		gameInstance.Player1.Move = &move
		return nil
	} else if gameInstance.Player2 != nil && gameInstance.Player2.Name == playerName {
		gameInstance.Player2.Move = &move
		return nil
	} else {
		return errors.New("player not found")
	}
}

func HaveAllPlayersSelectedMove() (bool, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if gameInstance.Player1 == nil || gameInstance.Player2 == nil {
		return false, errors.New("not all players have been set")
	}

	return gameInstance.Player1.Move != nil && gameInstance.Player2.Move != nil, nil
}

func GetRound() int {
	mutex.Lock()
	defer mutex.Unlock()
	return gameInstance.Round
}

func IncrementRound() {
	mutex.Lock()
	defer mutex.Unlock()
	gameInstance.Round++
}

func GetPlayerFromConn(conn *websocket.Conn) (*types.Player, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if gameInstance.Player1 != nil && gameInstance.Player1.Conn == conn {
		return gameInstance.Player1, nil
	} else if gameInstance.Player2 != nil && gameInstance.Player2.Conn == conn {
		return gameInstance.Player2, nil
	} else {
		return nil, errors.New("player not found for given connection")
	}
}

func GetOpponent(playerName string) (*types.Player, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if gameInstance.Player1 != nil && gameInstance.Player1.Name == playerName {
		if gameInstance.Player2 != nil {
			return gameInstance.Player2, nil
		} else {
			return nil, errors.New("opponent not found")
		}
	} else if gameInstance.Player2 != nil && gameInstance.Player2.Name == playerName {
		if gameInstance.Player1 != nil {
			return gameInstance.Player1, nil
		} else {
			return nil, errors.New("opponent not found")
		}
	} else {
		return nil, errors.New("player not found")
	}
}

func RemovePlayer(playerName string) error {
	mutex.Lock()
	defer mutex.Unlock()

	if gameInstance.Player1 != nil && gameInstance.Player1.Name == playerName {
		gameInstance.Player1 = nil
		return nil
	} else if gameInstance.Player2 != nil && gameInstance.Player2.Name == playerName {
		gameInstance.Player2 = nil
		return nil
	} else {
		return errors.New("player not found")
	}
}
