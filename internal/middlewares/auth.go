package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/utilities-payment-api-gateway/internal/grpc/auth"
	"github.com/iskanye/utilities-payment-api-gateway/internal/lib/jwt"
	"github.com/iskanye/utilities-payment-utils/pkg/logger"
)

const (
	prefix    = "Bearer "
	prefixLen = 7
)

func AuthMiddleware(a auth.Auth, log *slog.Logger, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Auth.Validation"

		log := log.With(
			slog.String("op", op),
		)

		log.Info("attempting to validate user")

		token := c.Request.Header.Get("Authorization")

		if !strings.HasPrefix(token, prefix) {
			log.Error("failed to get user token")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token = token[prefixLen:]

		userID, isAdmin, err := jwt.ValidateToken(token, secret)
		if err != nil {
			log.Error("failed to validate user", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err": err.Error(),
			})
			return
		}

		log.Info("validated successfully")
		c.Request.Header.Add("UserID", fmt.Sprint(userID))
		if isAdmin {
			c.Request.Header.Add("Admin", "1")
		}

		c.Next()
	}
}

func AdminMiddleware(a auth.Auth, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Auth.IsAdmin"

		log := log.With(
			slog.String("op", op),
		)

		log.Info("attempting to check users permissions")

		isAdmin := c.Request.Header.Get("Admin") != ""

		if !isAdmin {
			log.Warn("user is not admin")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err": "user is not admin",
			})
		}

		c.Next()
	}
}
