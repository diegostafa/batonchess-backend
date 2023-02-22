package main

import (
	"fmt"
)

const (
	addr      = "localhost"
	http_port = "2023"
	tcp_port  = "2024"
)

func main() {
	go BatonChessHttp(fmt.Sprintf("%s:%s", addr, http_port))
	go BatonChessTcp(fmt.Sprintf("%s:%s", addr, tcp_port))
	select {}
}
