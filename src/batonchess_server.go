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
	router.GET("/users/:id", getUserById)

	// game
	router.GET("/activeGames", getActiveGames)
	router.POST("/createGame", createGame)
}

func createUser(c *gin.Context) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return
	}

	if !InsertUser(uuid.String()) {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusAccepted, nil)

	user, err := GetUser(uuid.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func getUserById(c *gin.Context) {
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

func getActiveGames(c *gin.Context) {
}

func createGame(c *gin.Context) {
}
