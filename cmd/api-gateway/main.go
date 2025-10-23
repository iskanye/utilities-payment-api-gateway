package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iskanye/utilities-payment-api-gateway/internal/app"
	"github.com/iskanye/utilities-payment-api-gateway/internal/config"
	pkgConfig "github.com/iskanye/utilities-payment-api-gateway/pkg/config"
)

func main() {
	pkgConfig.MustLoad(func(t *config.Config) {})
	app := app.New(gin.Default())
	app.MustRun()
}
