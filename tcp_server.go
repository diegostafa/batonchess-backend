package main

import (
	"encoding/json"

	"github.com/firstrow/tcp_server"
)

func BatonChessTcp(addr string, be *BatonchessEngine) {
	server := tcp_server.New(addr)

	server.OnNewClient(onNewClientClosure(be))
	server.OnClientConnectionClosed(onCloseClientClosure(be))
	server.OnNewMessage(onMessageClientClosure(be))

	server.Listen()
}

func onNewClientClosure(be *BatonchessEngine) func(*tcp_server.Client) {
	return func(c *tcp_server.Client) {
		println("NEW CLIENT")
	}
}

func onCloseClientClosure(be *BatonchessEngine) func(*tcp_server.Client, error) {
	return func(c *tcp_server.Client, err error) {
		println("CLOSE CLIENT")
	}
}

func onMessageClientClosure(be *BatonchessEngine) func(*tcp_server.Client, string) {
	return func(c *tcp_server.Client, message string) {
		println("NEW MESSAGE")

		var action BatonchessTcpAction
		err := json.Unmarshal([]byte(message), &action)
		if err != nil {
			panic(err)
		}

		switch action.ActionType {
		case JoinGameAction:
			joinGameHandler(be, c, action.ActionBody)
		case MakeMoveAction:
			makeMoveHandler(be, c, action.ActionBody)
		}

	}
}

func joinGameHandler(be *BatonchessEngine, c *tcp_server.Client, reqData interface{}) {
	var (
		reqMap    (map[string]interface{})
		joinReq   JoinGameRequest
		player    UserPlayer
		gameId    GameId
		gameState *GameState
	)
	reqMap = reqData.(map[string]interface{})

	joinReq = JoinGameRequest{
		GameId:      int(reqMap["gameId"].(float64)),
		UserId:      reqMap["userId"].(string),
		UserName:    reqMap["userName"].(string),
		PlayAsWhite: reqMap["playAsWhite"].(bool),
	}

	player = UserPlayer{
		Id:             joinReq.UserId,
		Name:           joinReq.UserName,
		PlayingAsWhite: joinReq.PlayAsWhite,
	}

	playerTcp := UserPlayerTcp{
		ConnId: "",
		Player: player,
	}

	gameId = GameId{
		Id: joinReq.GameId,
	}

	gameState = be.joinGame(&playerTcp, &gameId)
	jsonBytes, _ := json.Marshal(gameState)
	c.Send(string(jsonBytes))
}

func makeMoveHandler(be *BatonchessEngine, c *tcp_server.Client, req interface{}) {}
