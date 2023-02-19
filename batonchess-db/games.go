package batonchessDb

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func CreateGame(gp *GameProps) (*GameInfo, error) {
	var (
		gameInfo GameInfo
		err      error
	)

	err = queryNone(`
		INSERT INTO
			games(creator_id,max_players,fen)
		VALUES
			(?,?,"fenn")`,
		gp.CreatorId, gp.MaxPlayers)

	if err != nil {
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
		FROM users_in_games RIGHT JOIN (
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

func GetPlayerInGames(user *UserId) (*UserId, error) {

	return nil, nil
}

func JoinGame(joinRequest *JoinGameRequest) error {
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
		SELECT
			COUNT(1)
		FROM
			users_in_games
		WHERE
			game_id = ? AND
			user_id = ?`,
		joinRequest.GameId,
		joinRequest.UserId)

	if err != nil {
		println(err.Error())
		return err
	}

	if exists == 0 {
		return queryNone(`
		INSERT INTO
			users_in_games(game_id, user_id, is_playing, play_as_white)
		VALUES
			(?,?,true,?)`,
			joinRequest.GameId,
			joinRequest.UserId,
			joinRequest.PlayAsWhite,
		)
	} else {
		return queryNone(`
		UPDATE
			users_in_games
		SET
			is_playing = true,
			play_as_white = ?
		WHERE
			game_id = ? AND user_id = ?`,
			joinRequest.PlayAsWhite,
			joinRequest.GameId,
			joinRequest.UserId,
		)
	}
}

func LeaveGame(leaveRequest *UsersInGamesId) error {
	return queryNone(`
		UPDATE
			users_in_games
		SET
			is_playing = false
		WHERE
			game_id = ? AND user_id = ?`,
		leaveRequest.GameId,
		leaveRequest.UserId,
	)
}

func GetGameInfo(gameId *GameId) (*GameInfo, error) {
	var (
		gameInfo GameInfo
		err      error
	)

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
			FROM games
			WHERE g_id = ?
		)
		AS game_info
		ON game_id = game_info.g_id
		GROUP BY game_info.g_id
		`, gameId.Id)

	if err != nil {
		return nil, err
	}

	return &gameInfo, nil
}

func GetGamePlayers(gameId *GameId) ([]Player, error) {
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
		FROM users_in_games RIGHT JOIN (
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

func GetGameState(gameId *GameId) (*GameState, error) {
	var (
		gameState *GameState
		gameInfo  *gameInfo
		players   []Player
		exists    int
		err       error
	)

	gameInfo, err = GetGameInfo(gameId)
	if err != nil {
		return nil, err
	}

	// get players
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
		FROM users_in_games RIGHT JOIN (
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
		println(err.Error())
		return nil, err
	}

	return gameState, nil
}
