package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func BatonChessHttp(addr string, be *BatonchessEngine) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	bindEndpoints(router, be)
	router.Run(addr)
}

func bindEndpoints(router *gin.Engine, be *BatonchessEngine) {
	router.GET("/createUser", createUser)
	router.POST("/updateUserName", updateUserName)
	router.POST("/isValidUser", isValidUser)
	router.GET("/getActiveGames", getActiveGamesClosure(be))
	router.POST("/createGame", createGameClosure(be))
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
	var updateInfo UpdateUsernameRequest

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
		println(err.Error())
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

func getActiveGamesClosure(be *BatonchessEngine) func(*gin.Context) {
	return func(c *gin.Context) {
		gameInfos, err := GetActiveGames()

		if err != nil {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		gamesInMemory := make([]GameInfo, 0, len(gameInfos))
		for _, value := range gameInfos {
			if _, ok := be.games[value.GameId]; ok {
				currPlayers := len(be.games[value.GameId].players)
				value.CurrentPlayers = currPlayers
				gamesInMemory = append(gamesInMemory, value)
			}
		}

		c.JSON(http.StatusOK, gamesInMemory)
	}
}

func createGameClosure(be *BatonchessEngine) func(*gin.Context) {
	return func(c *gin.Context) {
		var (
			gp       CreateGameRequest
			gameInfo *GameInfo
			err      error
		)

		if err := c.BindJSON(&gp); err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		gameInfo, err = CreateGame(&gp)

		if err != nil {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		gid := &GameId{Id: gameInfo.GameId}
		be.createGame(gid)
		gameInfo.CurrentPlayers = len(be.games[gid.Id].players)

		c.JSON(http.StatusCreated, gameInfo)
	}
}
