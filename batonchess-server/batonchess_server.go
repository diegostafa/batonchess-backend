package batonchessServer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BatonChessServer struct {
	router *gin.Engine
}

func NewBatonChessServer() *BatonChessServer {
	router := gin.Default()
	bindEndpoints(router)

	return &BatonChessServer{
		router: router,
	}
}

func (bc *BatonChessServer) listenOn(addr string) {
	bc.router.Run(addr)
}

func bindEndpoints(router *gin.Engine) {
	// users
	router.GET("/createUser", createUser)
	router.POST("/updateUserName", updateUserName)
	router.POST("/isValidUser", isValidUser)

	// game
	router.GET("/getActiveGames", getActiveGames)
	router.POST("/createGame", createGame)
	router.POST("/joinGame", joinGame)
	router.POST("/leaveGame", leaveGame)
}

// --- USER

func createUser(c *gin.Context) {
	uuid, err := uuid.NewRandom()

	if err != nil {
		return
	}

	user := User{
		Id:   uuid.String(),
		Name: "anon",
	}
	if err := CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusCreated, user)
}

func updateUserName(c *gin.Context) {
	var updateInfo UserNameUpdateRequest

	if err := c.BindJSON(&updateInfo); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if err := UpdateUserName(&updateInfo); err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusAccepted, nil)
}

func isValidUser(c *gin.Context) {
	var u UserId

	if err := c.BindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	res, err := GetUser(&u)

	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	if res == nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// --- GAME

func getActiveGames(c *gin.Context) {
	games, err := GetActiveGames()

	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, games)
}

func createGame(c *gin.Context) {
	var (
		gp   GameProps
		game *GameInfo
		err  error
	)

	if err := c.BindJSON(&gp); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	game, err = CreateGame(&gp)

	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusCreated, game)
}

func joinGame(c *gin.Context) {
	var (
		joinRequest JoinGameRequest
		gameState   *GameState
		err         error
	)

	if err := c.BindJSON(&joinRequest); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if err := JoinGame(&joinRequest); err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	gameState, err = GetGameState(&GameId{Id: joinRequest.GameId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusAccepted, gameState)
}

func leaveGame(c *gin.Context) {
	var leaveRequest UsersInGamesId

	if err := c.BindJSON(&leaveRequest); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if err := LeaveGame(&leaveRequest); err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, nil)
}
