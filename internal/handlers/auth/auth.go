package auth

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/utilities-payment-api-gateway/internal/grpc/auth"
	"github.com/iskanye/utilities-payment/pkg/logger"
)

func LoginHandler(a auth.Auth, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Auth.Login"

		email := c.Query("email")
		password := c.Query("password")

		log := log.With(
			slog.String("op", op),
			slog.String("email", email),
		)

		log.Info("attempting to login user")

		token, err := a.Login(c, email, password)
		if err != nil {
			log.Error("failed to login user", logger.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		log.Info("success")
		// c.SetCookie("token", token, cookieTTL, "/", host, false, true)
		// c.SetCookie("user_id", fmt.Sprint(userId), cookieTTL, "/", host, false, true)
		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}

func RegisterHandler(a auth.Auth, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Auth.Register"

		email := c.PostForm("email")
		password := c.PostForm("password")

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

		log.Info("successfully registered",
			slog.Int64("user_id", id),
		)

		c.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	}
}
