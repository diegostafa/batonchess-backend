batonchess_db = ./db/batonchess.db
setup_db = ./db/setup_db.sql

clean:
	rm $(batonchess_db)
	sqlite3 $(batonchess_db) < $(setup_db)

run:
	go run *.go
