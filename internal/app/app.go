package app

import (
	"log/slog"
	"net"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/utilities-payment-api-gateway/internal/config"
	grpcAuth "github.com/iskanye/utilities-payment-api-gateway/internal/grpc/auth"
	grpcbilling "github.com/iskanye/utilities-payment-api-gateway/internal/grpc/billing"
	authHandlers "github.com/iskanye/utilities-payment-api-gateway/internal/handlers/auth"
	billingHandlers "github.com/iskanye/utilities-payment-api-gateway/internal/handlers/billing"
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

	billing, err := grpcbilling.New(cfg.Billing.Host, cfg.Billing.Port)
	if err != nil {
		panic(err)
	}

	authMiddleware := middlewares.AuthMiddleware(&auth, log)

	// AUTH SERVICE
	login := authHandlers.LoginHandler(cfg, &auth, log)
	register := authHandlers.RegisterHandler(&auth, log)

	// BILLING SERVICE
	addBill := billingHandlers.AddBillHandler(&billing, log)

	engine.GET("/login", login)
	engine.GET("/register", register)

	// Auth required
	engine.Use(authMiddleware)
	{
		engine.GET("/addbill", addBill)
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
