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
	CreatorId               string `json:"id"`
	MaxPlayersPerSide       int    `json:"maxPlayersPerSide"`
	MinutesPerSide          int    `json:"minutesPerSide"`
	SecondsIncrementPerMove int    `json:"secondsIncrementPerMove"`
}

type GameState struct {
	Id    string `json:"id"`
	State int    `json:"State"`
}
