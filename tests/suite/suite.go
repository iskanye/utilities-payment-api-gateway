package suite

import (
	"context"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/utilities-payment-api-gateway/internal/app"
	"github.com/iskanye/utilities-payment-api-gateway/internal/config"
	pkgConfig "github.com/iskanye/utilities-payment/pkg/config"
)

type Suite struct {
	Cfg *config.Config

	e     *gin.Engine
	ctx   context.Context
	token string
}

func NewTest(t *testing.T) *Suite {
	gin.SetMode(gin.TestMode)
	t.Helper()
	t.Parallel()

	cfg := pkgConfig.MustLoadPath[config.Config](configPath())

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	eng := gin.New()
	app.New(eng, discardLogger(), cfg)

	return &Suite{
		e:   eng,
		Cfg: cfg,
		ctx: ctx,
	}
}

func NewBench(b *testing.B) *Suite {
	gin.SetMode(gin.TestMode)
	b.Helper()

	cfg := pkgConfig.MustLoadPath[config.Config](configPath())

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.Timeout)

	b.Cleanup(func() {
		b.Helper()
		cancelCtx()
	})

	eng := gin.New()
	app.New(eng, discardLogger(), cfg)

	return &Suite{
		e:   eng,
		Cfg: cfg,
		ctx: ctx,
	}
}

func configPath() string {
	const key = "CONFIG_PATH"

	if v := os.Getenv(key); v != "" {
		return v
	}

	return "../config/tests.yaml"
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
