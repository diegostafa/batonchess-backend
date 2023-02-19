package batonchess

import (
	"fmt"
	"os"
)

func main() {
	port := os.Args[1]
	NewBatonChessServer().listenOn(fmt.Sprintf("localhost:%s", port))
}
