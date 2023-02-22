package main

import (
	"encoding/json"

	"github.com/firstrow/tcp_server"
)

func BatonChessTcp(addr string) {
	server := tcp_server.New(addr)

	server.OnNewClient(onNewClientHandler)
	server.OnNewMessage(onNewMessageHandler)
	server.OnClientConnectionClosed(onClientConnectionClosedHandler)

	server.Listen()
}

func onNewClientHandler(c *tcp_server.Client) {
	println("NEW CLIENT")
}

func onNewMessageHandler(c *tcp_server.Client, message string) {
	println("NEW MESSAGE")

	var action BatonchessTcpAction
	err := json.Unmarshal([]byte(message), &action)
	if err != nil {
		panic(err)
	}

	switch action.ActionType {
	case JoinGameAction:
		println("JOIN GAME ACTION")
	case MakeMoveAction:
		println("MAKE MOVE ACTION")

	}

}

func onClientConnectionClosedHandler(c *tcp_server.Client, err error) {
	println("CLIENT CLOSE CONNECTION")
}
