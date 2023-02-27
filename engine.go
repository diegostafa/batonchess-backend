package main

import (
	"github.com/notnil/chess"
)

type BatonchessGame struct {
	maxPlayers  int
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

func (be *BatonchessEngine) createGame(gameInfo *GameInfo) {
	be.games[gameInfo.GameId] = &BatonchessGame{
		maxPlayers:  gameInfo.MaxPlayers,
		whiteQueue:  make([]UserPlayer, 0),
		blackQueue:  make([]UserPlayer, 0),
		fen:         INITIAL_FEN,
		isWhiteTurn: true}
}

func (be *BatonchessEngine) joinGame(player UserPlayer, gid *GameId) bool {
	game := be.games[gid.Id]

	if be.getCurrentPlayers(gid) == game.maxPlayers {
		return false
	} else if player.PlayingAsWhite {
		game.whiteQueue = append(game.whiteQueue, player)
	} else {
		game.blackQueue = append(game.blackQueue, player)
	}

	return true
}

func (be *BatonchessEngine) leaveGame(ug *UserInGame) {
	game := be.games[ug.GameId]

	for i, p := range game.whiteQueue {
		if p.Id == ug.UserId && i >= 0 && i < len(game.whiteQueue) {
			game.whiteQueue = append(game.whiteQueue[:i], game.whiteQueue[i+1:]...)
		}
	}

	for i, p := range game.blackQueue {
		if p.Id == ug.UserId && i >= 0 && i < len(game.blackQueue) {
			game.blackQueue = append(game.blackQueue[:i], game.blackQueue[i+1:]...)
		}
	}
}

func (be *BatonchessEngine) updateGame(updateReq *UpdateFenRequest) {
	game := be.games[updateReq.GameId]
	gameState := be.getGameState(&GameId{updateReq.GameId})

	if gameState.WaitingForPlayers {
		return
	}

	if gameState.UserToPlay.Id != updateReq.UserId {
		return
	}

	if !isValidNewPosition(game.fen, updateReq.NewFen) {
		return
	}

	if gameState.UserToPlay.PlayingAsWhite {
		game.whiteQueue = append(game.whiteQueue[1:], game.whiteQueue[0])
	} else {
		game.whiteQueue = append(game.whiteQueue[1:], game.whiteQueue[0])
	}

	game.fen = updateReq.NewFen
	game.isWhiteTurn = !game.isWhiteTurn
}

func (be *BatonchessEngine) getGameState(gid *GameId) *GameState {
	game := be.games[gid.Id]
	gameState := &GameState{}
	gameState.Fen = game.fen
	gameState.WhiteQueue = game.whiteQueue
	gameState.BlackQueue = game.blackQueue

	gameState.Outcome, gameState.Method = getChessboardState(game.fen)

	if len(game.whiteQueue) == 0 || len(game.blackQueue) == 0 {
		gameState.WaitingForPlayers = true
	} else if game.isWhiteTurn {
		gameState.UserToPlay = game.whiteQueue[0]
	} else {
		gameState.UserToPlay = game.blackQueue[0]
	}

	return gameState
}

func isValidNewPosition(currFenStr string, nextFenStr string) bool {
	currFen, err := chess.FEN(currFenStr)
	if err != nil {
		return false
	}

	_, err = chess.FEN(nextFenStr)
	if err != nil {
		return false
	}

	currGame := *chess.NewGame(currFen)
	validMoves := currGame.ValidMoves()

	for _, move := range validMoves {
		if currGame.Position().Update(move).String() == nextFenStr {
			return true
		}
	}

	return false
}

func getChessboardState(fenStr string) (string, int) {
	fen, _ := chess.FEN(fenStr)
	game := chess.NewGame(fen)
	outcome := game.Outcome().String()
	method := game.Method()
	return outcome, int(method)
}

func (be *BatonchessEngine) getCurrentPlayers(gid *GameId) int {
	game := be.games[gid.Id]
	return len(game.whiteQueue) + len(game.blackQueue)
}
