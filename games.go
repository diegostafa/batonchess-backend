package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func CreateGame(gp *CreateGameRequest) (*GameInfo, error) {
	var (
		gameInfo GameInfo
		err      error
	)

	err = queryNone(`
		INSERT INTO games(creator_id,max_players,fen)
		VALUES (?,?,"fenn")`,
		gp.CreatorId, gp.MaxPlayers)

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
				&gameInfo.GameStatus,
				&gameInfo.CreatedAt,
				&gameInfo.MaxPlayers)
		},
		`
		SELECT
			g_id,
			u_id,
			u_name,
			g_state,
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
					&g.GameStatus,
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
			g_state,
			created_at,
			max_players
		FROM games JOIN users
		ON games.creator_id = users.u_id
		WHERE g_state = "NORMAL"
		ORDER BY created_at DESC`)

	if err != nil {
		println(err.Error())
		return nil, err
	}

	return gameInfos, nil
}
