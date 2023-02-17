package main

import (
	"fmt"
	"os"
)

func main() {
	port := os.Args[1]

	setupDb()
	NewBatonChessServer().listenOn(fmt.Sprintf("localhost:%s", port))
}
