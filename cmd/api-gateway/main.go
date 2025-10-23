package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iskanye/utilities-payment-api-gateway/internal/app"
)

func main() {
	app := app.New(gin.Default())
	app.MustRun()
}
