batonchess_db = ./db/batonchess.db
server_port = 2023

clean:
	rm $(batonchess_db)

run:
	go run src/*.go ${server_port}