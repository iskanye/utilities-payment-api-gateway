package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/utilities-payment-api-gateway/internal/grpc/auth"
	"github.com/iskanye/utilities-payment-api-gateway/internal/lib/jwt"
	"github.com/iskanye/utilities-payment-utils/pkg/logger"
)

// POST /users/login
func LoginHandler(a auth.Auth, log *slog.Logger, tokenSaver jwt.TokenSaver, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Auth.Login"

		email := c.PostForm("email")
		password := c.PostForm("password")

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

		// Токен ОБЯЗАН быть валидным иначе сервис аутентификации неправильно работает
		payload, err := jwt.ValidateToken(token, secret)
		if err != nil {
			log.Error("failed to validate token", logger.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": err.Error(),
			})
			return
		}

		tokenSaver.Set(token, payload)

		log.Info("success")
		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}

// POST /users/register
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

// GET /admin/users
func GetUsersHandler(a auth.Auth, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Auth.GetUsers"

		log := log.With(
			slog.String("op", op),
		)

		log.Info("attempting to get users list")

		users, err := a.GetUsers(c)
		if err != nil {
			log.Error("failed to fetch users", logger.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		log.Info("users list fetched successfully")

		c.JSON(http.StatusOK, users)
	}
}
