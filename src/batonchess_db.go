package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var (
	sqliteVersion = "sqlite3"
	batonchessDb  = "./db/batonchess.db"
	setupDbScript = "./db/setup_db.sql"
)

// --- SETUP

func setupDb() error {
	data, err := os.ReadFile(setupDbScript)
	if err != nil {
		return err
	}

	if _, err := os.Stat(batonchessDb); os.IsNotExist(err) {
		return queryNone(string(data))
	}
	return nil
}

// --- DB

func queryNone(queryString string, queryArgs ...any) error {
	db, err := sql.Open(sqliteVersion, batonchessDb)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	stmt, err := tx.Prepare(queryString)
	if err != nil {
		println(err.Error())
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(queryArgs...)
	if err != nil {
		return fmt.Errorf("error executing query: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func queryOne(queryString string, queryArgs ...any) (*sql.Row, error) {
	db, err := sql.Open(sqliteVersion, batonchessDb)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	stmt, err := tx.Prepare(queryString)
	if err != nil {
		println(err.Error())
		return nil, fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(queryArgs...)

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return row, nil
}

func queryMany(queryString string, queryArgs ...any) (*sql.Rows, error) {
	db, err := sql.Open(sqliteVersion, batonchessDb)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	stmt, err := tx.Prepare(queryString)
	if err != nil {
		println(err.Error())
		return nil, fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	row, err := stmt.Query(queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return row, nil
}

// --- USER

func GetUser(uuid string) (*User, error) {
	row, err := queryOne(`SELECT * FROM users WHERE u_id = ?`, uuid)
	if err != nil {
		return nil, err
	}

	if row == nil {
		return nil, nil
	}

	user := User{}
	err = row.Scan(&user.Id, &user.Name)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func InsertUser(uuid string) bool {
	return nil == queryNone(`INSERT INTO users(u_id, u_name) VALUES(?, DEFAULT)`, uuid)
}

func UpdateUserName(updateInfo *UserNameUpdateRequest) bool {
	return nil == queryNone(
		`UPDATE users SET u_name = ? WHERE u_id = ?`,
		updateInfo.NewUsername, updateInfo.UserId,
	)
}

// --- GAME

func CreateGame(user *User, gp *GameProps) bool {
	return nil == queryNone(`
		INSERT INTO games(g_id,fen,is_active,max_players,minutes_per_side, seconds_increment)
		VALUES(DEFAULT, DEFAULT, DEFAULT,?,?,?)`,
		gp.MaxPlayersPerSide, gp.MinutesPerSide, gp.SecondsIncrementPerMove)
}
