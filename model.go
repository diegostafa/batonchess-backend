package main

const (
	JOIN_GAME_ACTION  = "JOIN_GAME"
	UPDATE_FEN_ACTION = "UPDATE_FEN"
	REFUSED_ACTION    = "REFUSED"
	INITIAL_FEN       = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

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
	JoinedAt       int64  `json:"joinedAt"`
}

type UserInGame struct {
	UserId string `json:"userId"`
	GameId int    `json:"gameId"`
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
	WhiteQueue        []UserPlayer `json:"whiteQueue"`
	BlackQueue        []UserPlayer `json:"blackQueue"`
	UserToPlay        UserPlayer   `json:"userToPlay"`
	WaitingForPlayers bool         `json:"waitingForPlayers"`
	Outcome           string       `json:"outcome"`
	Method            int          `json:"method"`
}

type CreateGameRequest struct {
	CreatorId  string `json:"creatorId"`
	MaxPlayers int    `json:"maxPlayers"`
}

type GameInfo struct {
	GameId         int    `json:"gameId"`
	CreatorId      string `json:"creatorId"`
	CreatorName    string `json:"creatorName"`
	CreatedAt      int64  `json:"createdAt"`
	MaxPlayers     int    `json:"maxPlayers"`
	CurrentPlayers int    `json:"currentPlayers"`
}

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
