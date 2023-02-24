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
		var (
			action     BatonchessTcpAction
			actionBody []byte
			err        error
		)
		err = json.Unmarshal([]byte(message), &action)
		if err != nil {
			panic(err)
		}

		actionBody, err = json.Marshal(action.ActionBody)
		if err != nil {
			panic(err)
		}

		switch action.ActionType {
		case JOIN_GAME_ACTION:
			joinGameHandler(be, c, actionBody)
		case UPDATE_FEN_ACTION:
			updateFenHandler(be, c, actionBody)
		case LEAVE_GAME_ACTION:
			leaveGameHandler(be, c, actionBody)

		}
	}
}

func joinGameHandler(be *BatonchessEngine, c *tcp_server.Client, jsonReq []byte) {
	var joinReq JoinGameRequest
	err := json.Unmarshal(jsonReq, &joinReq)
	if err != nil {
		panic(err)
	}

	playerTcp := UserPlayerTcp{
		ConnId: c.Conn().RemoteAddr().String(),
		Player: UserPlayer{
			Id:             joinReq.UserId,
			Name:           joinReq.UserName,
			PlayingAsWhite: joinReq.PlayAsWhite,
		},
	}

	gameState := be.joinGame(playerTcp, &GameId{Id: joinReq.GameId})
	gameStateJson, _ := json.Marshal(gameState)
	c.Send(string(gameStateJson))
}

func updateFenHandler(be *BatonchessEngine, c *tcp_server.Client, jsonReq []byte) {
	var updateReq UpdateFenRequest
	err := json.Unmarshal(jsonReq, &updateReq)
	if err != nil {
		panic(err)
	}

	gameState, success := be.updateFen(&updateReq)
	if !success {
		c.Send(DISCARDED_MOVE_MESSAGE)
		return
	}

	gameStateJson, _ := json.Marshal(*gameState)
	c.Send(string(gameStateJson))
}

func leaveGameHandler(be *BatonchessEngine, c *tcp_server.Client, jsonReq []byte) {
	be.leaveGame(&UserPlayerTcp{}, &GameId{})
}
