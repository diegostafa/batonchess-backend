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

func queryOne(scanner func(*sql.Row) error, queryString string, queryArgs ...any) error {
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

	row := stmt.QueryRow(queryArgs...)
	scanned := scanner(row)
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return scanned
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
	var user User
	err := queryOne(
		func(row *sql.Row) error {
			println("GOT USERSSS")
			return row.Scan(&user.Id, &user.Name)
		},
		`SELECT * FROM users WHERE u_id = ?`, u.Id)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *User) bool {
	return nil == queryNone(`INSERT INTO users(u_id, u_name) VALUES(?, ?)`, user.Id, user.Name)
}

func UpdateUserName(updateInfo *UserNameUpdateRequest) bool {
	return nil == queryNone(
		`UPDATE users SET u_name = ? WHERE u_id = ?`,
		updateInfo.NewUsername, updateInfo.Id,
	)
}

// --- GAME

func CreateGame(gp *GameProps) (*GameInfo, error) {
	err := queryNone(`
		INSERT INTO games(creator_id,max_players,fen)
		VALUES(?,?,"fenn")`,
		gp.CreatorId, gp.MaxPlayers)

	if err != nil {
		return nil, err
	}

	var gameInfo GameInfo
	err = queryOne(
		func(row *sql.Row) error {
			return row.Scan(
				&gameInfo.GameId,
				&gameInfo.CreatorName,
				&gameInfo.GameStatus,
				&gameInfo.CreatedAt,
				&gameInfo.MaxPlayers,
				&gameInfo.CurrentPlayers)
		},
		`
		SELECT
			g_id,
			u_name,
			g_state,
			created_at,
			max_players,
			COUNT(users_in_games.user_id) AS current_players
		FROM users_in_games RIGHT JOIN (
			SELECT *
			FROM games JOIN users
			ON games.creator_id = users.u_id
			WHERE creator_id = ?
			ORDER BY created_at DESC
			LIMIT 1
		) AS game_info
		ON game_id = game_info.g_id
		GROUP BY game_info.g_id
		`, gp.CreatorId)

	if err != nil {
		println(err.Error())
		return nil, err
	}

	return &gameInfo, nil
}

func GetUserActiveGame() {}

func GetActiveGames() ([]GameInfo, error) {
	rows, err := queryMany(`
		SELECT *
		FROM games
		WHERE g_state = ?`, "NORMAL")
	if err != nil {
		return nil, err
	}

	if rows == nil {
		return nil, nil
	}

	var games []GameInfo
	for rows.Next() {
		var g GameInfo
		err := rows.Scan(&g.GameId, &g.CreatorName, &g.GameStatus, &g.MaxPlayers, &g.CurrentPlayers)
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
