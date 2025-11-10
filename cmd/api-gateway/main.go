package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/utilities-payment-api-gateway/internal/app"
	"github.com/iskanye/utilities-payment-api-gateway/internal/config"
	pkgConfig "github.com/iskanye/utilities-payment-utils/pkg/config"
	"github.com/iskanye/utilities-payment-utils/pkg/logger"
)

func main() {
	cfg := pkgConfig.MustLoad[config.Config]()
	cfg.MustLoadSecret()

	log := setupPrettySlog()
	app := app.New(gin.Default(), log, cfg)

	go func() {
		app.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	log.Info("Gracefully stopped")
}

func setupPrettySlog() *slog.Logger {
	opts := logger.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
