package main

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserNameUpdateRequest struct {
	UserId      string `json:"userId"`
	NewUsername string `json:"newUsername"`
}

type GameProps struct {
	CreatorId         string `json:"creatorId"`
	MaxPlayersPerSide int    `json:"maxPlayersPerSide"`
	PlayAsWhite       bool   `json:"playAsWhite"`
	// SecondsIncrementPerMove int    `json:"secondsIncrementPerMove"`
	// MinutesPerSide          int    `json:"minutesPerSide"`
}

type GameState struct {
	Id    string `json:"id"`
	State int    `json:"State"`
}
