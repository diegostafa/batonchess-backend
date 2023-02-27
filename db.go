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

func queryMany(scanner func(*sql.Rows) error, queryString string, queryArgs ...any) error {
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

	rows, err := stmt.Query(queryArgs...)
	if err != nil {
		return fmt.Errorf("error executing query: %v", err)
	}

	scanned := scanner(rows)

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return scanned
}

func GetUser(u *UserId) (*User, error) {
	var (
		user User
		err  error
	)
	err = queryOne(
		func(row *sql.Row) error {
			return row.Scan(&user.Id, &user.Name)
		}, `
		SELECT *
		FROM users
		WHERE u_id = ?`,
		u.Id)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateUser(user *User) error {
	return queryNone(`
		INSERT INTO users(u_id, u_name)
		VALUES (?, ?)`,
		user.Id, user.Name)
}

func UpdateUserName(updateInfo *UpdateUsernameRequest) error {
	return queryNone(`
		UPDATE users
		SET u_name = ?
		WHERE u_id = ?`,
		updateInfo.NewName, updateInfo.Id,
	)
}

func CreateGame(gp *CreateGameRequest) (*GameInfo, error) {
	var (
		gameInfo GameInfo
		err      error
	)

	err = queryNone(`
		INSERT INTO games(creator_id,max_players,fen)
		VALUES (?,?,?)`,
		gp.CreatorId, gp.MaxPlayers, INITIAL_FEN)

	if err != nil {
		println(err.Error())
		return nil, err
	}

	err = queryOne(
		func(row *sql.Row) error {
			return row.Scan(
				&gameInfo.GameId,
				&gameInfo.CreatorId,
				&gameInfo.CreatorName,
				&gameInfo.CreatedAt,
				&gameInfo.MaxPlayers)
		},
		`
		SELECT
			g_id,
			u_id,
			u_name,
			created_at,
			max_players
		FROM games JOIN users
		ON games.creator_id = users.u_id
		WHERE creator_id = ?
		ORDER BY created_at DESC
		LIMIT 1`,
		gp.CreatorId)

	if err != nil {
		println(err.Error())
		return nil, err
	}

	return &gameInfo, nil
}

func GetActiveGames() ([]GameInfo, error) {
	var (
		gameInfos []GameInfo
		g         GameInfo
		err       error
	)

	err = queryMany(
		func(rows *sql.Rows) error {
			for rows.Next() {
				err := rows.Scan(
					&g.GameId,
					&g.CreatorId,
					&g.CreatorName,
					&g.CreatedAt,
					&g.MaxPlayers)
				if err != nil {
					return err
				}
				gameInfos = append(gameInfos, g)
			}
			return nil
		},
		`
		SELECT
			g_id,
			u_id,
			u_name,
			created_at,
			max_players
		FROM games JOIN users
		ON games.creator_id = users.u_id
		WHERE outcome = "*"
		ORDER BY created_at DESC`)

	if err != nil {
		println(err.Error())
		return nil, err
	}

	return gameInfos, nil
}

func UpdateGameState(gid *GameId, gameState *GameState) error {
	return queryNone(`
		UPDATE games
		SET outcome = ?, method = ?
		WHERE g_id = ?`,
		gameState.Outcome, gameState.Method, gid.Id,
	)
}
