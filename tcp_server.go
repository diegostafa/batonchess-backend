package main

import (
	"encoding/json"

	"github.com/firstrow/tcp_server"
)

var (
	conns map[*tcp_server.Client]UserInGame
)

func BatonChessTcp(addr string, be *BatonchessEngine) {
	server := tcp_server.New(addr)
	conns = make(map[*tcp_server.Client]UserInGame)

	server.OnNewClient(onNewClient)
	server.OnClientConnectionClosed(onCloseClientClosure(be))
	server.OnNewMessage(onMessageClientClosure(be))
	server.Listen()
}

func broadcastGameState(gid *GameId, gameState *GameState) {
	println("BROADCASTING")
	var clients []*tcp_server.Client
	for k, v := range conns {
		if v.GameId == gid.Id {
			clients = append(clients, k)
		}
	}

	for i, c := range clients {
		println("BROADCASTING: ", i)

		gameStateJson, _ := json.Marshal(*gameState)
		c.Send(string(gameStateJson))
	}

}

func onNewClient(c *tcp_server.Client) {
	println("NEW CLIENT")
}

func onCloseClientClosure(be *BatonchessEngine) func(*tcp_server.Client, error) {
	return func(c *tcp_server.Client, err error) {
		ug := conns[c]
		gameState := be.leaveGame(&ug)
		if gameState == nil {
			return
		}

		broadcastGameState(&GameId{ug.GameId}, gameState)
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
			println("JOIN GAME")
			joinGameHandler(be, c, actionBody)
		case UPDATE_FEN_ACTION:
			println("UPDATE FEN")
			updateFenHandler(be, c, actionBody)
		}
	}
}

func joinGameHandler(be *BatonchessEngine, c *tcp_server.Client, jsonReq []byte) {
	var joinReq JoinGameRequest
	err := json.Unmarshal(jsonReq, &joinReq)
	if err != nil {
		panic(err)
	}

	player := UserPlayer{
		Id:             joinReq.UserId,
		Name:           joinReq.UserName,
		PlayingAsWhite: joinReq.PlayAsWhite,
	}

	gameState := be.joinGame(player, &GameId{Id: joinReq.GameId})
	if gameState == nil {
		println("JOIN REFUSED")
		c.Send(REFUSED_ACTION)
		return
	}

	conns[c] = UserInGame{UserId: joinReq.UserId, GameId: joinReq.GameId}
	broadcastGameState(&GameId{joinReq.GameId}, gameState)
}

func updateFenHandler(be *BatonchessEngine, c *tcp_server.Client, jsonReq []byte) {
	var updateReq UpdateFenRequest
	err := json.Unmarshal(jsonReq, &updateReq)
	if err != nil {
		panic(err)
	}

	gameState := be.updateFen(&updateReq)
	if gameState == nil {
		c.Send(REFUSED_ACTION)
		return
	}

	broadcastGameState(&GameId{updateReq.GameId}, gameState)
}
