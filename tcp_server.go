package main

import (
	"encoding/json"

	"github.com/firstrow/tcp_server"
)

var (
	userOfClient map[*tcp_server.Client]*UserInGame
)

func BatonChessTcp(addr string, be *BatonchessEngine) {
	server := tcp_server.New(addr)
	userOfClient = make(map[*tcp_server.Client]*UserInGame)

	server.OnNewMessage(onMessageClientClosure(be))
	server.OnClientConnectionClosed(onCloseClientClosure(be))
	server.Listen()
}

func onMessageClientClosure(be *BatonchessEngine) func(*tcp_server.Client, string) {
	return func(c *tcp_server.Client, message string) {
		println("ON NEW MESSAGE")
		var (
			action     BatonchessTcpAction
			actionBody []byte
			err        error
		)
		err = json.Unmarshal([]byte(message), &action)
		if err != nil {
			println(err.Error())
			return
		}

		actionBody, err = json.Marshal(action.ActionBody)
		if err != nil {
			println(err.Error())
			return
		}

		switch action.ActionType {
		case JOIN_GAME_ACTION:
			joinGameHandler(be, c, actionBody)
		case UPDATE_FEN_ACTION:
			updateFenHandler(be, c, actionBody)
		}

	}
}

func onCloseClientClosure(be *BatonchessEngine) func(*tcp_server.Client, error) {
	return func(c *tcp_server.Client, err error) {
		println("ON CLOSE CLIENT")
		if err != nil {
			println("CLOSE ERR:", err.Error())
		}

		ug := userOfClient[c]
		if ug == nil {
			broadcastGameState(nil, nil)
			println("NO USER OF CLIENT")
			return
		}

		gameState := be.leaveGame(ug)
		if gameState == nil {
			println("NO GAME STATE")
			return
		}

		delete(userOfClient, c)
		broadcastGameState(&GameId{ug.GameId}, gameState)
	}
}

func broadcastGameState(gid *GameId, gameState *GameState) {
	println("BROADCASTING NEW GAME STATE")
	if gid == nil || gameState == nil {
		return
	}
	var clients []*tcp_server.Client

	for k, v := range userOfClient {
		if v.GameId == gid.Id {
			clients = append(clients, k)
		}
	}

	for _, c := range clients {
		gameStateJson, _ := json.Marshal(*gameState)
		c.Send(string(gameStateJson))
	}
}

func joinGameHandler(be *BatonchessEngine, c *tcp_server.Client, jsonReq []byte) {
	var joinReq JoinGameRequest
	err := json.Unmarshal(jsonReq, &joinReq)
	if err != nil {
		println(err.Error())
		return
	}

	player := UserPlayer{
		Id:             joinReq.UserId,
		Name:           joinReq.UserName,
		PlayingAsWhite: joinReq.PlayAsWhite,
	}

	gameState := be.joinGame(player, &GameId{Id: joinReq.GameId})
	if gameState == nil {
		c.Send(REFUSED_ACTION)
		return
	}

	userOfClient[c] = &UserInGame{UserId: joinReq.UserId, GameId: joinReq.GameId}
	broadcastGameState(&GameId{joinReq.GameId}, gameState)
}

func updateFenHandler(be *BatonchessEngine, c *tcp_server.Client, jsonReq []byte) {
	var updateReq UpdateFenRequest
	err := json.Unmarshal(jsonReq, &updateReq)
	if err != nil {
		println(err.Error())
		return
	}

	gameState := be.updateFen(&updateReq)
	if gameState == nil {
		c.Send(REFUSED_ACTION)
		return
	}

	broadcastGameState(&GameId{updateReq.GameId}, gameState)
}
