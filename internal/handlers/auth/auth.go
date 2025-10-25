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

		log := log.With(
			slog.String("op", op),
		)

		email := c.Query("email")
		password := c.Query("password")

		log.Info("attempting to login user",
			slog.String("email", email),
		)

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

		log := log.With(
			slog.String("op", op),
		)

		email := c.Query("email")
		password := c.Query("password")

		log.Info("attempting to register user",
			slog.String("email", email),
		)

		id, err := a.Register(c, email, password)
		if err != nil {
			log.Error("failed to register user", logger.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"id":  0,
				"err": err.Error(),
			})
			return
		}

		log.Info("success")

		c.JSON(http.StatusOK, gin.H{
			"id":  id,
			"err": "nil",
		})
	}
}
