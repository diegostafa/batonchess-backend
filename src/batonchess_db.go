package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var (
	sqliteVersion = "sqlite3"
	batonchessDb  = "./db/batonchess.db"
)

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

func GetUser(u *UserId) (*User, error) {
	row, err := queryOne(`SELECT * FROM users WHERE u_id = ?`, u.Id)
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

func InsertUser(user *User) bool {
	return nil == queryNone(`INSERT INTO users(u_id, u_name) VALUES(?, ?)`, user.Id, user.Name)
}

func UpdateUserName(updateInfo *UserNameUpdateRequest) bool {
	return nil == queryNone(
		`UPDATE users SET u_name = ? WHERE u_id = ?`,
		updateInfo.NewUsername, updateInfo.UserId,
	)
}

// --- GAME

func CreateGame(gp *GameProps) (int, error) {
	err := queryNone(`
		INSERT INTO games(g_id,creator_id,fen,state,max_players)
		VALUES(DEFAULT,?,DEFAULT,DEFAULT,?)`,
		gp.CreatorId, gp.MaxPlayersPerSide)

	if err != nil {
		return -1, err
	}

	return 0, nil
}

func GetUserActiveGame() {}

func GetActiveGames() ([]Game, error) {
	rows, err := queryMany(`SELECT * FROM games WHERE g_state = ?`, "NORMAL")
	if err != nil {
		return nil, err
	}

	if rows == nil {
		return nil, nil
	}

	var games []Game
	for rows.Next() {
		var g Game
		err := rows.Scan(&g.Id, &g.Fen, &g.MaxPlayers, &g.Status)
		if err != nil {
			return nil, err
		}
		games = append(games, g)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return games, nil
}
