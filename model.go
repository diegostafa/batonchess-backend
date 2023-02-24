package main

const (
	JOIN_GAME_ACTION       = "JOIN_GAME"
	LEAVE_GAME_ACTION      = "LEAVE_GAME"
	UPDATE_FEN_ACTION      = "UPDATE_FEN"
	DISCARDED_MOVE_MESSAGE = "DISCARDED"
	INITIAL_FEN            = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

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

type UserPlayerTcp struct {
	ConnId string     `json:"connId"`
	Player UserPlayer `json:"player"`
}

func (upt *UserPlayerTcp) toUserPlayer() *UserPlayer {
	return &upt.Player
}

type UpdateUsernameRequest struct {
	Id      string `json:"id"`
	NewName string `json:"newName"`
}

type GameId struct {
	Id int `json:"id"`
}

type GameState struct {
	Fen               string       `json:"fen"`
	Players           []UserPlayer `json:"players"`
	UserToPlay        UserPlayer   `json:"userToPlay"`
	WaitingForPlayers bool         `json:"waitingForPlayers"`
}

type CreateGameRequest struct {
	CreatorId  string `json:"creatorId"`
	MaxPlayers int    `json:"maxPlayers"`
}

type GameInfo struct {
	GameId         int    `json:"gameId"`
	CreatorName    string `json:"creatorName"`
	GameStatus     string `json:"gameStatus"`
	CreatedAt      int    `json:"createdAt"`
	MaxPlayers     int    `json:"maxPlayers"`
	CurrentPlayers int    `json:"currentPlayers"`
}

// --- TCP SERVER

type BatonchessTcpAction struct {
	ActionType string      `json:"actionType"`
	ActionBody interface{} `json:"actionBody"`
}

type JoinGameRequest struct {
	GameId      int    `json:"gameId"`
	UserId      string `json:"userId"`
	UserName    string `json:"userName"`
	PlayAsWhite bool   `json:"playAsWhite"`
}

type UpdateFenRequest struct {
	GameId int    `json:"gameId"`
	UserId string `json:"userId"`
	NewFen string `json:"newFen"`
}
