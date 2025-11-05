package app

import (
	"log/slog"
	"net"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/utilities-payment-api-gateway/internal/config"
	"github.com/iskanye/utilities-payment-api-gateway/internal/gin/handlers"
	grpcAuth "github.com/iskanye/utilities-payment-api-gateway/internal/grpc/auth"
	grpcBilling "github.com/iskanye/utilities-payment-api-gateway/internal/grpc/billing"
	grpcPayment "github.com/iskanye/utilities-payment-api-gateway/internal/grpc/payment"
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

	billing, err := grpcBilling.New(cfg.Billing.Host, cfg.Billing.Port)
	if err != nil {
		panic(err)
	}

	payment, err := grpcPayment.New(cfg.Payment.Host, cfg.Payment.Port)
	if err != nil {
		panic(err)
	}

	// MIDDLEWARES
	authMiddleware := middlewares.AuthMiddleware(&auth, log, cfg.AuthSecret)
	adminsMiddleware := middlewares.AdminMiddleware(&auth, log)

	// AUTH SERVICE
	login := handlers.LoginHandler(&auth, log)
	register := handlers.RegisterHandler(&auth, log)

	// BILLING SERVICE
	addBill := handlers.AddBillHandler(&billing, log)
	getBill := handlers.GetBillHandler(&billing, log)
	getBills := handlers.GetBillsHandler(&billing, log)

	// PAYMENT SERVICE
	payBill := handlers.PayBillHandler(&payment, &billing, log)

	engine.POST("/users/login", login)
	engine.POST("/users/register", register)

	// Auth required
	authorized := engine.Group("/", authMiddleware)
	{
		admins := authorized.Group("/admin", adminsMiddleware)
		{
			admins.POST("/bills", addBill)
		}
		authorized.GET("/bills", getBills)
		authorized.GET("/bills/:id", getBill)
		authorized.POST("/bills/pay", payBill)
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
