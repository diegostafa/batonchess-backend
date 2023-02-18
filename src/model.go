package main

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserId struct {
	Id string `json:"id"`
}

type UserNameUpdateRequest struct {
	UserId      string `json:"userId"`
	NewUsername string `json:"newUsername"`
}

type GameProps struct {
	CreatorId         string `json:"creatorId"`
	MaxPlayersPerSide int    `json:"maxPlayersPerSide"`
	PlayAsWhite       bool   `json:"playAsWhite"`
}

type Game struct {
	Id         string `json:"id"`
	Fen        string `json:"fen"`
	MaxPlayers int    `json:"maxPlayers"`
	Status     string `json:"status"`
}
