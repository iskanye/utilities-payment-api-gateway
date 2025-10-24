package app

import (
	"log/slog"
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/utilities-payment-api-gateway/internal/config"
	grpcAuth "github.com/iskanye/utilities-payment-api-gateway/internal/grpc/auth"
	authHandlers "github.com/iskanye/utilities-payment-api-gateway/internal/handlers/auth"
	"github.com/iskanye/utilities-payment-api-gateway/internal/middlewares"
)

type App struct {
	e   *gin.Engine
	cfg *config.Config
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

	authMiddleware := middlewares.AuthMiddleware(&auth, log)

	// AUTH SERVICE
	login := authHandlers.LoginHandler(cfg, &auth, log)
	register := authHandlers.RegisterHandler(&auth, log)

	engine.GET("/login", login)
	engine.GET("/register", register)

	// Auth required
	engine.Use(authMiddleware)
	{
		engine.GET("/bills", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"msg": "your bills here",
			})
		})
	}

	return &App{
		e:   engine,
		cfg: cfg,
	}
}

func (a *App) MustRun() {
	if err := a.e.Run(net.JoinHostPort(a.cfg.Host, strconv.Itoa(a.cfg.Port))); err != nil {
		panic(err)
	}
}
