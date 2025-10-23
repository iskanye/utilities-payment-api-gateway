package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type App struct {
	e *gin.Engine
}

func New(
	engine *gin.Engine,
) *App {
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	return &App{e: engine}
}

func (a *App) MustRun() {
	if err := a.e.Run(); err != nil {
		panic(err)
	}
}
