package main

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserId struct {
	Id string `json:"id"`
}

type GameId struct {
	Id int `json:"id"`
}

type UserNameUpdateRequest struct {
	Id          string `json:"id"`
	NewUsername string `json:"newUsername"`
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

type GameState struct {
	Id         string `json:"id"`
	Fen        string `json:"fen"`
	MaxPlayers int    `json:"maxPlayers"`
	Status     string `json:"status"`
}
