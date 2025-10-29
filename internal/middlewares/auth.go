package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/utilities-payment-api-gateway/internal/grpc/auth"
	"github.com/iskanye/utilities-payment-api-gateway/internal/lib/jwt"
	"github.com/iskanye/utilities-payment/pkg/logger"
)

const prefix = "Bearer "

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

		token = token[len(prefix):]

		userID, err := jwt.ValidateToken(token, secret)
		if err != nil {
			log.Error("failed to validate user", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err": err.Error(),
			})
			return
		}

		log.Error("validated successfully")
		c.Request.Header.Add("UserID", fmt.Sprint(userID))
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

		idStr := c.Request.Header.Get("UserID")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Error("cant convert id to int64", logger.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{
				"err": "cant convert id to int64",
			})
			return
		}

		isAdmin, err := a.IsAdmin(c, id)

		if !isAdmin {
			log.Warn("user is not admin", logger.Err(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err": "user is not admin",
			})
		}

		c.Next()
	}
}
