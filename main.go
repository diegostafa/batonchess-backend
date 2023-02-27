package main

const (
	http_addr = "192.168.1.100:2023"
	tcp_addr  = "192.168.1.100:2024"
)

func main() {
	be := NewBatonchessEngine()
	go BatonChessHttp(http_addr, be)
	go BatonChessTcp(tcp_addr, be)
	select {}
}
