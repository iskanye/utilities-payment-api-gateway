package auth

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/utilities-payment-api-gateway/internal/grpc/auth"
	"github.com/iskanye/utilities-payment/pkg/logger"
)

func LoginHandler(a auth.Auth, log *slog.Logger) func(*gin.Context) {
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

		token, err := a.Login(c, email, password)
		if err != nil {
			log.Error("failed to login user", logger.Err(err))
			c.JSON(http.StatusBadRequest, map[string]string{
				"token": "",
				"err":   err.Error(),
			})
			return
		}

		log.Info("success")

		c.JSON(http.StatusOK, map[string]string{
			"token": token,
			"err":   "null",
		})
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
			c.JSON(http.StatusBadRequest, map[string]any{
				"id":  0,
				"err": err.Error(),
			})
			return
		}

		log.Info("success")

		c.JSON(http.StatusOK, map[string]any{
			"id":  id,
			"err": "nil",
		})
	}
}
