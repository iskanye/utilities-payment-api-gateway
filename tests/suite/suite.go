package suite

import (
	"context"
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

	e   *gin.Engine
	ctx context.Context
}

func NewTest(t *testing.T) *Suite {
	t.Helper()
	t.Parallel()

	cfg := pkgConfig.MustLoadPath[config.Config](configPath())

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	eng := gin.Default()
	app.New(eng, slog.Default(), cfg)

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
