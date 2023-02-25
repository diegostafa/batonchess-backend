package main

type BatonchessGame struct {
	players     []UserPlayer
	whiteQueue  []UserPlayer
	blackQueue  []UserPlayer
	fen         string
	isWhiteTurn bool
}

type BatonchessEngine struct {
	games map[int]*BatonchessGame
}

func NewBatonchessEngine() *BatonchessEngine {
	return &BatonchessEngine{
		games: make(map[int]*BatonchessGame),
	}
}

func (be *BatonchessEngine) createGame(gid *GameId) {
	be.games[gid.Id] = &BatonchessGame{
		players:     make([]UserPlayer, 0),
		whiteQueue:  make([]UserPlayer, 0),
		blackQueue:  make([]UserPlayer, 0),
		fen:         INITIAL_FEN,
		isWhiteTurn: true}
}

func (be *BatonchessEngine) joinGame(player UserPlayer, gid *GameId) *GameState {
	game := be.games[gid.Id]
	game.players = append(game.players, player)

	if player.PlayingAsWhite {
		game.whiteQueue = append(game.whiteQueue, player)
	} else {
		game.blackQueue = append(game.blackQueue, player)
	}

	return be.getGameState(gid)
}

func (be *BatonchessEngine) leaveGame(ug *UserInGame) *GameState {
	game := be.games[ug.GameId]

	for i, p := range game.players {
		if p.Id == ug.UserId && i >= 0 && i < len(game.players) {
			game.players = append(game.players[:i], game.players[i+1:]...)
		}
	}

	return be.getGameState(&GameId{ug.GameId})
}

func (be *BatonchessEngine) updateFen(updateReq *UpdateFenRequest) *GameState {
	game := be.games[updateReq.GameId]
	gameState := be.getGameState(&GameId{updateReq.GameId})

	if gameState.WaitingForPlayers {
		return nil
	}

	if gameState.UserToPlay.Id != updateReq.UserId {
		return nil
	}

	if gameState.UserToPlay.PlayingAsWhite {
		game.whiteQueue = append(game.whiteQueue[1:], game.whiteQueue[0])
	} else {
		game.whiteQueue = append(game.whiteQueue[1:], game.whiteQueue[0])
	}

	game.fen = updateReq.NewFen
	gameState.Fen = game.fen
	game.isWhiteTurn = !game.isWhiteTurn

	return gameState
}

func (be *BatonchessEngine) getGameState(gid *GameId) *GameState {
	game := be.games[gid.Id]
	gameState := &GameState{}
	gameState.Players = game.players
	gameState.Fen = game.fen

	if len(game.whiteQueue) == 0 || len(game.blackQueue) == 0 {
		gameState.WaitingForPlayers = true
	} else if game.isWhiteTurn {
		gameState.UserToPlay = game.whiteQueue[0]
	} else {
		gameState.UserToPlay = game.blackQueue[0]
	}

	return gameState
}
