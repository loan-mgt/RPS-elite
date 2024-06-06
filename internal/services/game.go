package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"rcp/elite/internal/senders"
	"rcp/elite/internal/types"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	playerPoll         = map[string]types.Player{}
	playerPollMutex    = &sync.Mutex{}
	gamePoll           = map[string]types.Game{}
	gamePollMutex      = &sync.Mutex{}
	playerGameMap      = map[string]string{}
	playerGameMapMutex = &sync.Mutex{}
	objective          = 3
)

func isPlayerInGame(playerName string) (bool, string) {
	playerGameMapMutex.Lock()
	defer playerGameMapMutex.Unlock()
	gId, ok := playerGameMap[playerName]

	return ok, gId
}

func isPlayerInPool(playerName string) bool {
	playerPollMutex.Lock()
	defer playerPollMutex.Unlock()
	_, ok := playerPoll[playerName]

	return ok
}

func JoinPlayerPoll(player types.Player) error {
	if is, _ := isPlayerInGame(player.Name); is {
		return errors.New("Player is in game")
	}

	if isPlayerInPool(player.Name) {
		return nil
	}

	playerPoll[player.Name] = player
	return nil
}

func StartPollMonitor() {
	for {
		time.Sleep(3 * time.Second)
		playerPollMutex.Lock()
		for k, v := range playerPoll {
			err := sendPing(v)
			if err != nil {
				delete(playerPoll, k)
			}
		}
		playerPollMutex.Unlock()
	}
}

func sendPing(player types.Player) error {
	return player.Conn.WriteMessage(websocket.PingMessage, []byte(""))
}

func SearchForGameToCreate() {
	for {
		time.Sleep(3 * time.Second)
		playerPollMutex.Lock()
		var opponent *types.Player
		opponent = nil
		for k, v := range playerPoll {
			if opponent == nil {
				opponent = &v
			} else {
				createGame(*opponent, v)
				delete(playerPoll, k)
				delete(playerPoll, opponent.Name)
				opponent = nil
			}
		}
		playerPollMutex.Unlock()
	}

}

func createGame(player1, player2 types.Player) {

	gameId := generateGameId(player1.Name, player2.Name, time.Now().Unix())

	gamePollMutex.Lock()

	gamePoll[gameId] = types.Game{
		Players: map[string]types.Player{
			player1.Name: player1,
			player2.Name: player2,
		},
		Round: 1,
	}

	gamePollMutex.Unlock()

	playerGameMapMutex.Lock()

	playerGameMap[player1.Name] = gameId
	playerGameMap[player2.Name] = gameId

	playerGameMapMutex.Unlock()

	go mainGameLoop(gameId)

}

func generateGameId(player1Name, player2Name string, unixTime int64) string {
	combinedInput := fmt.Sprintf("%s%s%d", player1Name, player2Name, unixTime)

	hash := sha256.New()

	hash.Write([]byte(combinedInput))

	hashBytes := hash.Sum(nil)

	hashString := hex.EncodeToString(hashBytes)

	return hashString
}

func IsPlayerInGame(playerName string) bool {
	playerGameMapMutex.Lock()
	defer playerGameMapMutex.Unlock()
	_, ok := playerGameMap[playerName]
	return ok
}

func GetPlayerStatus(playerName string) (*types.Player, error) {
	playerGameMapMutex.Lock()
	v, ok := playerGameMap[playerName]
	playerGameMapMutex.Unlock()

	if !ok {
		return nil, errors.New("player not in game")
	}

	gamePollMutex.Lock()
	g, ok := gamePoll[v]
	gamePollMutex.Unlock()

	if !ok {
		return nil, errors.New("game not found")
	}

	p, ok := g.Players[playerName]

	if !ok {
		return nil, errors.New("player not found in game")
	}

	return &p, nil
}

func SetPlayerMove(playerName string, move string) error {

	playerGameMapMutex.Lock()
	v, ok := playerGameMap[playerName]
	playerGameMapMutex.Unlock()

	if !ok {
		return errors.New("player not in game")
	}

	gamePollMutex.Lock()
	g, ok := gamePoll[v]
	gamePollMutex.Unlock()

	if !ok {
		return errors.New("game not found")
	}

	p, ok := g.Players[playerName]

	if !ok {
		return errors.New("player not found in game")
	}

	p.Move = move

	g.Players[playerName] = p

	gamePollMutex.Lock()
	gamePoll[v] = g
	gamePollMutex.Unlock()

	return nil
}

func getGameFromId(gameId string) (types.Game, bool) {
	gamePollMutex.Lock()
	defer gamePollMutex.Unlock()

	v, ok := gamePoll[gameId]
	return v, ok
}

func mainGameLoop(gameId string) {

	setupPlayersScreen(gameId)

	//handle init for player game start in 3s
	time.Sleep(3 * time.Second)

	//game start
	// you have 5s to select a move
	for still2player(gameId) || noPlayerHasWin(gameId) {
		time.Sleep(5 * time.Second)

		_, _ = incrementWinner(gameId)

		incrementRound(gameId)

		break
	}

	g, ok := getGameFromId(gameId)

	if !ok {
		return
	}

	playerGameMapMutex.Lock()
	for _, p := range g.Players {
		delete(playerGameMap, p.Name)
	}
	playerGameMapMutex.Unlock()

	gamePollMutex.Lock()
	delete(gamePoll, gameId)
	gamePollMutex.Unlock()

}

func setupPlayersScreen(gameId string) {
	g, ok := getGameFromId(gameId)

	if !ok {
		return
	}

	for _, p := range g.Players {
		senders.SetGameHome(p.Conn, &p, getOpponent(g.Players, p), "Game will start in 3s", "empty-3s")
	}
}

func getOpponent(players map[string]types.Player, player types.Player) *types.Player {

	for _, p := range players {
		if p.Name != player.Name {
			return &p
		}
	}
	return nil
}

func still2player(gameId string) bool {
	g, ok := getGameFromId(gameId)

	if !ok {
		return false
	}

	return len(g.Players) != 2
}

func noPlayerHasWin(gameId string) bool {
	g, ok := getGameFromId(gameId)

	if !ok {
		return false
	}

	for _, p := range g.Players {
		if p.Score >= 3 {
			return true
		}
	}

	return false

}

func incrementRound(gameId string) {
	g, ok := getGameFromId(gameId)

	if !ok {
		return
	}
	g.Round += 1

	saveGame(gameId, g)
}

func saveGame(gameId string, game types.Game) {
	gamePollMutex.Lock()
	defer gamePollMutex.Unlock()

	gamePoll[gameId] = game
}

func incrementWinner(gameId string) (string, bool) {
	g, ok := getGameFromId(gameId)

	if !ok {
		return "", false
	}

	winner := ""
	winnerMove := ""

	for k, p := range g.Players {
		if isLeftWinning(p.Move, winnerMove) {
			winner = k
			winnerMove = p.Move
		}
	}

	if winner != "" {
		tmpPlayer := g.Players[winner]
		tmpPlayer.Score += 1

		g.Players[winner] = tmpPlayer
	}

	saveGame(gameId, g)

	return winner, winner == ""

}

func isLeftWinning(leftMove, rightMove string) bool {
	if leftMove == "paper" && rightMove == "rock" {
		return true
	}

	if leftMove == "rock" && rightMove == "scissors" {
		return true
	}

	if leftMove == "scissors" && rightMove == "paper" {
		return true
	}

	return false

}
