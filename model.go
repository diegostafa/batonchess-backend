package main

// --- HTTP SERVER

type UserId struct {
	Id string `json:"id"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserPlayer struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	PlayingAsWhite bool   `json:"playingAsWhite"`
}

type UpdateUsernameRequest struct {
	Id      string `json:"id"`
	NewName string `json:"newName"`
}

type GameId struct {
	Id int `json:"id"`
}

type GameState struct {
	Fen        string       `json:"fen"`
	UserIdTurn string       `json:"userIdTurn"`
	Players    []UserPlayer `json:"players"`
}

type CreateGameRequest struct {
	CreatorId  string `json:"creatorId"`
	MaxPlayers int    `json:"maxPlayers"`
}

type GameInfo struct {
	GameId         int    `json:"gameId"`
	CreatorName    string `json:"creatorName"`
	GameStatus     string `json:"gameStatus"`
	CreatedAt      string `json:"createdAt"`
	MaxPlayers     int    `json:"maxPlayers"`
	CurrentPlayers int    `json:"currentPlayers"`
}

// --- TCP SERVER

const (
	JoinGameAction string = "join_game"
	MakeMoveAction string = "make_move"
)

type BatonchessTcpAction struct {
	ActionType string      `json:"actionType"`
	ActionBody interface{} `json:"actionBody"`
}

type JoinGameTcpRequest struct {
	GameId      int    `json:"gameId"`
	UserId      string `json:"userid"`
	PlayAsWhite bool   `json:"playAsWhite"`
}

type MakeMoveTcpRequest struct {
	GameId      int    `json:"gameId"`
	UserId      string `json:"userid"`
	PlayAsWhite bool   `json:"playAsWhite"`
}
