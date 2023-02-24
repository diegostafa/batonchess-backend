package main

type BatonchessGame struct {
	playersTcp  []UserPlayerTcp
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
		playersTcp:  make([]UserPlayerTcp, 0),
		whiteQueue:  make([]UserPlayer, 0),
		blackQueue:  make([]UserPlayer, 0),
		fen:         INITIAL_FEN,
		isWhiteTurn: true}
}

func (be *BatonchessEngine) joinGame(userTcp UserPlayerTcp, gid *GameId) *GameState {
	game := be.games[gid.Id]
	game.playersTcp = append(game.playersTcp, userTcp)

	if userTcp.Player.PlayingAsWhite {
		game.whiteQueue = append(game.whiteQueue, userTcp.Player)
	} else {
		game.blackQueue = append(game.blackQueue, userTcp.Player)
	}

	return be.getGameState(gid)
}

func (be *BatonchessEngine) leaveGame(userTcp *UserPlayerTcp, gid *GameId) *GameState {
	game := be.games[gid.Id]
	for i, p := range game.playersTcp {
		if p.Player.Id == userTcp.Player.Id {
			game.playersTcp = append(game.playersTcp[:i], game.playersTcp[i+1:]...)
		}
	}

	return be.getGameState(gid)
}

func (be *BatonchessEngine) updateFen(updateReq *UpdateFenRequest) (*GameState, bool) {
	gameState := be.getGameState(&GameId{Id: updateReq.GameId})
	gameState.Fen = updateReq.NewFen
	be.games[updateReq.GameId].isWhiteTurn = !be.games[updateReq.GameId].isWhiteTurn
	return gameState, true
}

func (be *BatonchessEngine) getGameState(gid *GameId) *GameState {
	var (
		gameState *GameState
		game      *BatonchessGame
		players   []UserPlayer
	)

	game = be.games[gid.Id]
	gameState = &GameState{}

	for i := 0; i < len(game.playersTcp); i++ {
		players = append(players, *game.playersTcp[i].toUserPlayer())
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
