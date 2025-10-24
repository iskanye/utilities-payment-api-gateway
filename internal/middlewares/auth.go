package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/utilities-payment-api-gateway/internal/grpc/auth"
	"github.com/iskanye/utilities-payment/pkg/logger"
)

func AuthMiddleware(a auth.Auth, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Auth.Validation"

		log := log.With(
			slog.String("op", op),
		)

		log.Info("attempting to validate user")

		token, err := c.Cookie("token")
		if err != nil {
			log.Error("failed to get user token", logger.Err(err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		isValid, err := a.Validate(c, token)
		if err != nil {
			log.Error("failed to validate user", logger.Err(err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if isValid {
			log.Info("validated successfully")
			c.Next()
			return
		}

		log.Error("invalid token")
		c.Abort()
	}
}
