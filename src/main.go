package main

func main() {
	setupDb()
	server := NewBatonChessServer()
	server.listenOn("localhost:2023")
}
