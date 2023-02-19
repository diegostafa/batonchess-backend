package batonchess

type UserId struct {
	Id string `json:"id"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Player struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	PlayingAsWhite bool   `json:"playingAsWhite"`
}

type UserNameUpdateRequest struct {
	Id          string `json:"id"`
	NewUsername string `json:"newUsername"`
}

type GameId struct {
	Id int `json:"id"`
}

type GameState struct {
	GameInfo  GameInfo `json:"gameInfo"`
	Fen       string   `json:"fen"`
	Players   []Player `json:"players"`
	TurnQueue []UserId `json:"turnQueue"`
}

type GameProps struct {
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

type UsersInGamesId struct {
	GameId int    `json:"gameId"`
	UserId string `json:"userid"`
}

type JoinGameRequest struct {
	GameId      int    `json:"gameId"`
	UserId      string `json:"userid"`
	PlayAsWhite bool   `json:"playAsWhite"`
}
