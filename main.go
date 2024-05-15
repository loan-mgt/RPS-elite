package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Game struct {
	ID          string     `json:"id"`
	Players     [2]Player  `json:"players"`
	GameStatus  GameStatus `json:"gameStatus"`
	LastUpdated time.Time  `json:"lastUpdated"`
	Type        string     `json:"type"`
}

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

type GameStatus struct {
	PlayerScore   int    `json:"playerScore"`
	OpponentScore int    `json:"opponentScore"`
	GameMode      string `json:"gameMode"`
	Time          int    `json:"time"`
	Message       string `json:"message"`
	Status        string `json:"status"`
	Round         int    `json:"round"`
}

var games map[string]*Game
var gamesMutex sync.Mutex

func init() {
	games = make(map[string]*Game)
}

func wsHandlerTest(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Send message back to the client indicating successful join
	sendMessage(conn, Message{Type: "test", Success: true})

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Received message: %s type: %s\n", message, string(rune(messageType)))

		// Execute the template with any required data
		tmpl, err := template.ParseFiles("component/gameHome.html")
		if err != nil {
			log.Println("Error parsing template:", err)
			return
		}

		// Create a buffer to execute the template into
		var tplBuffer bytes.Buffer
		err = tmpl.Execute(&tplBuffer, nil)
		if err != nil {
			log.Println("Error executing template:", err)
			return
		}

		// Write the template content as a text message to the client
		err = conn.WriteMessage(websocket.TextMessage, tplBuffer.Bytes())
		if err != nil {
			log.Println("Error writing template to client:", err)
			return
		}

		log.Println(err)
	}

}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Read game ID from request URL
	vars := mux.Vars(r)
	gameID := vars["id"]

	// Get or create game
	game, err := getOrCreateGame(gameID)
	if err != nil {
		log.Println(err)
		return
	}

	// Add player to the game
	var moveData struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}

	err = conn.ReadJSON(&moveData)
	if err != nil {
		log.Println(err)
		return
	}

	playerID := fmt.Sprintf("%d", time.Now().UnixNano())
	player := Player{
		ID:    playerID,
		Name:  moveData.Name,
		Move:  "",
		Flag:  "default",
		Score: 0,
		Conn:  conn,
	}

	opponent := Player{}
	if game.Players[0].ID == "" {
		game.Players[0] = player
		opponent = game.Players[1]
	} else if game.Players[1].ID == "" {
		game.Players[1] = player
		opponent = game.Players[0]
	} else {
		log.Println("Game is full")
		return
	}

	// Send message back to the client indicating successful join
	sendMessage(conn, Message{Type: "join_ack", Success: true})

	// Send game details to player
	err = sendGameDetails(conn, game)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Received message: %s type: %s\n", message, string(rune(messageType)))

		game, err := getOrCreateGame(gameID)
		if err != nil {
			log.Println(err)
			return
		}

		if game.Players[0].ID == playerID {
			player = game.Players[0]
			opponent = game.Players[1]
		} else {
			player = game.Players[1]
			opponent = game.Players[0]
		}

		// Parse JSON message to get player's move
		var moveData struct {
			Type string `json:"type"`
			Move string `json:"move"`
		}
		if err := json.Unmarshal(message, &moveData); err != nil {
			log.Println(err)
			return
		}

		player.Move = moveData.Move

		// print current move
		fmt.Printf("Player: %s, Move: %s\n", player, moveData.Move)

		// check if every one has selected a move
		if player.Move == "" {

			game.GameStatus.Message = "Please make your move"
		} else if opponent.Move == "" {

			game.GameStatus.Message = fmt.Sprintf("Waiting for %s\n", opponent.Name)

		} else {

			winner, tie := determineWinner(player, opponent)

			responseMsg := ""

			if tie {
				responseMsg = "it's a tie"
			} else {
				if winner.ID == player.ID {
					player.Score++
				} else {
					opponent.Score++
				}
				responseMsg = fmt.Sprintf("%s won the round", winner.Name)
			}

			// Prepare game result
			game.GameStatus.Message = fmt.Sprintf("Player: %s -> %s, Opponent: %s -> %s, %s",
				player.Name,
				player.Move,
				opponent.Name,
				opponent.Move,
				responseMsg)

			// reseting for next round
			game.GameStatus.Round++
			player.Move = ""
			opponent.Move = ""

			fmt.Printf("Player: %s, Opponent: %s\n", player, opponent)

		}

		game.Players = [2]Player{player, opponent}
		game.Type = "game_update"

		fmt.Print("Game players: ", game.Players)

		saveGame(game)

		// Send game details to player type game_update
		err = broadcastGameDetails(game)
		if err != nil {
			log.Println(err)
			return
		}

	}
}

// Send game details to all players
func broadcastGameDetails(game *Game) error {
	for _, player := range game.Players {
		err := sendGameDetails(player.Conn, game)
		if err != nil {
			return err
		}
	}
	return nil
}

// Function to generate a random move for the server
func getRandomMove() string {
	moves := []string{"rock", "paper", "scissors"}
	randIndex := rand.Intn(len(moves))
	return moves[randIndex]
}

// Function to determine the winner and update the score
func determineWinner(player, opponent Player) (winner Player, tie bool) {
	playerMove := player.Move
	serverMove := opponent.Move

	if playerMove == serverMove {
		return Player{}, true
	} else if (playerMove == "rock" && serverMove == "scissors") ||
		(playerMove == "paper" && serverMove == "rock") ||
		(playerMove == "scissors" && serverMove == "paper") {
		return player, false
	} else {
		return opponent, false
	}
}

func getOrCreateGame(gameID string) (*Game, error) {
	gamesMutex.Lock()
	defer gamesMutex.Unlock()

	if game, ok := games[gameID]; ok {
		return game, nil
	}

	// Create a new game
	game := &Game{
		ID: gameID,
		Players: [2]Player{
			{ID: "", Name: "", Flag: "", Score: 0, Conn: nil, Move: ""},
			{ID: "", Name: "", Flag: "", Score: 0, Conn: nil, Move: ""},
		},
		GameStatus: GameStatus{
			PlayerScore:   0,
			OpponentScore: 0,
			GameMode:      "default",
			Time:          0,
			Message:       "",
			Status:        "waiting",
			Round:         0,
		},
		LastUpdated: time.Now(),
	}
	games[gameID] = game

	return game, nil
}

func saveGame(game *Game) error {
	gamesMutex.Lock()
	defer gamesMutex.Unlock()

	// Update the last updated time
	game.LastUpdated = time.Now()

	// Save the game to the games map
	games[game.ID] = game

	return nil
}

func sendGameDetails(conn *websocket.Conn, game *Game) error {
	gameJSON, err := json.Marshal(game)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, gameJSON)
	if err != nil {
		return err
	}

	return nil
}

func sendMessage(conn *websocket.Conn, message Message) error {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, messageJSON)
	if err != nil {
		return err
	}

	return nil
}

type Message struct {
	Type    string `json:"type"`
	Success bool   `json:"success,omitempty"`
}

func createGameHandler(w http.ResponseWriter, r *http.Request) {

	gameID := fmt.Sprintf("%d", time.Now().UnixNano())

	_, err := getOrCreateGame(gameID)

	if err != nil {
		w.Write([]byte("Failed to create game"))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write([]byte(gameID))
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ws/{id}", wsHandler)
	r.HandleFunc("/create", createGameHandler).Methods("POST")
	r.HandleFunc("/ws2", wsHandlerTest)

	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("").Handler(http.StripPrefix("", fs))

	http.Handle("/", r)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
