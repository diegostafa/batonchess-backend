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
		if p.Id == ug.UserId {
			game.players = append(game.players[:i], game.players[i+1:]...)
		}
	}

	return be.getGameState(&GameId{ug.GameId})
}

func (be *BatonchessEngine) updateFen(updateReq *UpdateFenRequest) *GameState {
	gameState := be.getGameState(&GameId{Id: updateReq.GameId})
	gameState.Fen = updateReq.NewFen
	be.games[updateReq.GameId].isWhiteTurn = !be.games[updateReq.GameId].isWhiteTurn
	return gameState
}

func (be *BatonchessEngine) getGameState(gid *GameId) *GameState {
	var (
		gameState *GameState
		game      *BatonchessGame
		players   []UserPlayer
	)

	game = be.games[gid.Id]
	gameState = &GameState{}

	for i := 0; i < len(game.players); i++ {
		players = append(players, game.players[i])
	}

	if len(game.whiteQueue) == 0 || len(game.blackQueue) == 0 {
		gameState.WaitingForPlayers = true
	} else if game.isWhiteTurn {
		gameState.UserToPlay = game.whiteQueue[0]
		game.whiteQueue = append(game.whiteQueue[1:], game.whiteQueue[0])
	} else {
		gameState.UserToPlay = game.blackQueue[0]
		game.blackQueue = append(game.blackQueue[1:], game.blackQueue[0])
	}

	gameState.Players = players
	gameState.Fen = game.fen
	return gameState
}
