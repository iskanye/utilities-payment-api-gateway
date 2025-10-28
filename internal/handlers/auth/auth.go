package auth

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/utilities-payment-api-gateway/internal/config"
	"github.com/iskanye/utilities-payment-api-gateway/internal/grpc/auth"
	"github.com/iskanye/utilities-payment/pkg/logger"
)

func LoginHandler(cfg *config.Config, a auth.Auth, log *slog.Logger) func(*gin.Context) {
	return func(c *gin.Context) {
		const op = "Auth.Login"

		email := c.Query("email")
		password := c.Query("password")

		log := log.With(
			slog.String("op", op),
			slog.String("email", email),
		)

		log.Info("attempting to login user")

		token, userId, err := a.Login(c, email, password)
		if err != nil {
			log.Error("failed to login user", logger.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		log.Info("success")
		c.SetCookie("token", token, int(cfg.CookieTTL.Seconds()), "/", cfg.Host, false, true)
		c.SetCookie("user_id", fmt.Sprint(userId), int(cfg.CookieTTL.Seconds()), "/", cfg.Host, false, true)
		c.JSON(http.StatusOK, nil)
	}
}

func RegisterHandler(a auth.Auth, log *slog.Logger) func(*gin.Context) {
	return func(c *gin.Context) {
		const op = "Auth.Register"

		email := c.Query("email")
		password := c.Query("password")

		log := log.With(
			slog.String("op", op),
			slog.String("email", email),
		)

		log.Info("attempting to register user")

		id, err := a.Register(c, email, password)
		if err != nil {
			log.Error("failed to register user", logger.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		log.Info("success")

		c.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	}
}
