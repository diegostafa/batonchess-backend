package main

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
	router.POST("/updateUserName/", updateUserName)
	router.POST("/isValidUser/", isValidUser)

	// game
	router.GET("/getActiveGames", getActiveGames)
	router.POST("/createGame/", createGame)
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
	if !InsertUser(&user) {
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

	if !UpdateUserName(&updateInfo) {
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
	var gp GameProps

	if err := c.BindJSON(&gp); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	game, err := CreateGame(&gp)

	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusCreated, game)
}
