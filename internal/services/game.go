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
	objective    = 3
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

func SetPlayerMove(playerName string, move *string) error {
	mutex.Lock()
	defer mutex.Unlock()

	if gameInstance.Player1 != nil && gameInstance.Player1.Name == playerName {
		gameInstance.Player1.Move = move
		return nil
	} else if gameInstance.Player2 != nil && gameInstance.Player2.Name == playerName {
		gameInstance.Player2.Move = move
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

func GetWinner() (player *types.Player, tie bool, err error) {
	mutex.Lock()
	defer mutex.Unlock()

	if gameInstance.Player1 == nil || gameInstance.Player2 == nil {
		return nil, false, errors.New("not all players have been set")
	}

	if gameInstance.Player1.Move == nil || gameInstance.Player2.Move == nil {
		return nil, false, errors.New("one or more player has not selected a move")
	}

	if gameInstance.Player1.Move == gameInstance.Player2.Move {
		return nil, true, nil
	}

	winingMove := getWinningMove(*gameInstance.Player1.Move, *gameInstance.Player2.Move)

	if winingMove == *gameInstance.Player1.Move {
		return gameInstance.Player1, false, nil
	}

	return gameInstance.Player2, false, nil

}

// waring tie are not handeled
func getWinningMove(move1 string, move2 string) string {
	if (move1 == "rock" || move2 == "rock") && (move1 == "paper" || move2 == "paper") {
		return "paper"
	} else if move1 == "rock" || move2 == "rock" {
		return "rock"
	} else {
		return "scissor"
	}
}

func IncrementPlayerScore(playerName string, amount int) error {
	mutex.Lock()
	defer mutex.Unlock()

	if playerName == gameInstance.Player1.Name {
		gameInstance.Player1.Score += amount
		return nil
	} else if playerName == gameInstance.Player2.Name {
		gameInstance.Player2.Score += amount
		return nil
	}
	return errors.New("unable to find player from name")

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

func IsGameFinish() (bool, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if gameInstance.Player1 == nil || gameInstance.Player2 == nil {
		return false, nil
	}

	return gameInstance.Player1.Score >= objective || gameInstance.Player2.Score >= objective, nil
}

func IsPlayerWinner(playerName string) (bool, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if gameInstance.Player1 != nil && gameInstance.Player1.Name == playerName {
		return gameInstance.Player1.Score >= objective, nil
	} else if gameInstance.Player2 != nil && gameInstance.Player2.Name == playerName {
		return gameInstance.Player2.Score >= objective, nil
	}
	return false, errors.New("unable to find player")

}
