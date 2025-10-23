package app

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/utilities-payment-api-gateway/internal/config"
	grpcAuth "github.com/iskanye/utilities-payment-api-gateway/internal/grpc/auth"
	authHandlers "github.com/iskanye/utilities-payment-api-gateway/internal/handlers/auth"
)

type App struct {
	e *gin.Engine
}

func New(
	engine *gin.Engine,
	log *slog.Logger,
	cfg *config.Config,
) *App {
	auth, err := grpcAuth.New(cfg.Auth.Host, cfg.Auth.Port)
	if err != nil {
		panic(err)
	}

	login := authHandlers.LoginHandler(&auth, log)
	register := authHandlers.RegisterHandler(&auth, log)

	engine.GET("/login", login)
	engine.GET("/register", register)

	return &App{e: engine}
}

func (a *App) MustRun() {
	if err := a.e.Run(); err != nil {
		panic(err)
	}
}
