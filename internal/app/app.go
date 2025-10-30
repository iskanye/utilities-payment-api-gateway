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

	// MIDDLEWARES
	authMiddleware := middlewares.AuthMiddleware(&auth, log, cfg.AuthSecret)
	adminsMiddleware := middlewares.AdminMiddleware(&auth, log)

	// AUTH SERVICE
	login := authHandlers.LoginHandler(&auth, log)
	register := authHandlers.RegisterHandler(&auth, log)

	// BILLING SERVICE
	addBill := billingHandlers.AddBillHandler(&billing, log)
	getBills := billingHandlers.GetBillsHandler(&billing, log)

	engine.GET("/user", login)
	engine.POST("/user", register)

	// Auth required
	authorized := engine.Group("/", authMiddleware)
	{
		admins := authorized.Group("/admin", adminsMiddleware)
		{
			admins.POST("/addbill", addBill)
		}
		authorized.GET("/bills", getBills)
	}

	return &App{
		e:   engine,
		cfg: cfg,
	}
}

func (a *App) MustRun() {
	if err := a.e.Run(address(a.cfg.Host, a.cfg.Port)); err != nil {
		panic(err)
	}
}

func address(host string, port int) string {
	return net.JoinHostPort(host, strconv.Itoa(port))
}
