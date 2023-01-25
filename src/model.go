package main

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserNameUpdateRequest struct {
	UserId      string `json:"userId"`
	NewUsername string `json:"newUsername`
}

type GameProps struct {
	UserId            string `json:"id"`
	MinsPerSide       int    `json:"minsPerSide"`
	SecsIncremeent    int    `json:"secsIncrement"`
	MaxPlayersPerSide int    `json:"maxPlayersPerSide"`
}

type GameState struct {
	Id    string `json:"id"`
	State int    `json:"State"`
}
