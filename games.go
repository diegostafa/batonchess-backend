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
		FROM users_in_games RIGHT JOIN
		(
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

func GetActiveGames() ([]GameInfo, error) {
	var (
		gameInfos []GameInfo
		err       error
	)

	err = queryMany(
		func(rows *sql.Rows) error {
			for rows.Next() {
				var g GameInfo
				err := rows.Scan(
					&g.GameId,
					&g.CreatorName,
					&g.GameStatus,
					&g.CreatedAt,
					&g.MaxPlayers,
					&g.CurrentPlayers)
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
			u_name,
			g_state,
			created_at,
			max_players,
			COUNT(users_in_games.user_id) AS current_players
		FROM users_in_games RIGHT JOIN
		(
			SELECT *
			FROM games JOIN users
			ON games.creator_id = users.u_id
			WHERE g_state = ?
			ORDER BY created_at DESC
		) AS game_info
		ON game_id = game_info.g_id
		GROUP BY game_info.g_id
		`, "NORMAL")

	if err != nil {
		return nil, err
	}

	return gameInfos, nil
}

func JoinGame(joinRequest *JoinGameRequest) error {
	println("JOINGAME")

	var (
		exists int
		err    error
	)

	err = queryOne(
		func(row *sql.Row) error {
			return row.Scan(
				&exists)
		},
		`
		SELECT COUNT(1)
		FROM users_in_games
		WHERE game_id = ? AND user_id = ?`,
		joinRequest.GameId,
		joinRequest.UserId)

	if err != nil {
		println(err.Error())
		return err
	}

	if exists == 0 {
		return queryNone(`
		INSERT INTO users_in_games(game_id, user_id, is_playing, play_as_white)
		VALUES (?,?,true,?)`,
			joinRequest.GameId,
			joinRequest.UserId,
			joinRequest.PlayAsWhite,
		)
	} else {
		return queryNone(`
		UPDATE users_in_games
		SET is_playing = true, play_as_white = ?
		WHERE game_id = ? AND user_id = ?`,
			joinRequest.PlayAsWhite,
			joinRequest.GameId,
			joinRequest.UserId,
		)
	}
}

func LeaveGame(leaveRequest *UsersInGamesId) error {
	return queryNone(`
		UPDATE users_in_games
		SET is_playing = false
		WHERE game_id = ? AND user_id = ?`,
		leaveRequest.GameId,
		leaveRequest.UserId,
	)
}

func GetGamePlayers(gameId *GameId) ([]UserPlayer, error) {
	println("GETGAMEPLAYERS")

	var (
		players []UserPlayer
		err     error
	)

	err = queryMany(
		func(rows *sql.Rows) error {
			for rows.Next() {
				var p UserPlayer
				err := rows.Scan(
					&p.Id,
					&p.Name,
					&p.PlayingAsWhite,
				)
				if err != nil {
					return err
				}
				players = append(players, p)
			}
			return nil
		},
		`
		SELECT users.u_id, users.u_name, g.play_as_white
		FROM users JOIN
		(
			SELECT user_id, play_as_white
			FROM users_in_games
			WHERE game_id = ?
		) AS g
		ON g.user_id = users.u_id
		`, gameId.Id)

	if err != nil {
		return nil, err
	}

	return players, nil
}

func GetGameState(gameId *GameId) (*GameState, error) {
	println("GETGAMESTATE")
	var (
		gameState *GameState
		players   []UserPlayer
		err       error
	)

	players, err = GetGamePlayers(gameId)
	if err != nil {
		return nil, err
	}

	gameState = &GameState{}
	gameState.Players = players
	gameState.Fen = "fen"
	gameState.UserIdTurn = "il diocan"

	return gameState, nil
}
