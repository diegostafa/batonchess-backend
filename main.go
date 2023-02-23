package main

const (
	http_addr = "localhost:2023"
	tcp_addr  = "localhost:2024"
)

func main() {
	be := NewBatonchessEngine()
	go BatonChessHttp(http_addr, &be)
	go BatonChessTcp(tcp_addr, &be)
	select {}
}
