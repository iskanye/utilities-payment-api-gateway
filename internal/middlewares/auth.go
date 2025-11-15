package middlewares

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iskanye/utilities-payment-api-gateway/internal/cache"
	"github.com/iskanye/utilities-payment-api-gateway/internal/grpc/auth"
	"github.com/iskanye/utilities-payment-api-gateway/internal/lib/jwt"
)

const (
	prefix    = "Bearer "
	prefixLen = 7
)

func AuthMiddleware(
	a auth.Auth,
	log *slog.Logger,
	tokenProvider jwt.TokenProvider,
	secret string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Auth.Validation"

		log := log.With(
			slog.String("op", op),
		)

		log.Info("attempting to validate user")

		token := c.Request.Header.Get("Authorization")

		if !strings.HasPrefix(token, prefix) {
			log.Error("failed to get user token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err": "no token provided",
			})
			return
		}

		token = token[prefixLen:]

		err := tokenProvider.Get(token)
		if err == cache.ErrCacheMiss {
			// Если токена нет в кеше значит он не заблочен
			payload, err := jwt.ValidateToken(token, secret)
			if err != nil {
				log.Warn("failed to validate")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"err": err.Error(),
				})
				return
			}

			log.Info("validated successfully")

			c.Set("Token", token)
			c.Set("UserID", payload.UserID)
			c.Set("IsAdmin", payload.IsAdmin)

			c.Next()
			return
		} else if err != nil {
			log.Warn("failed to access cache")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"err": err.Error(),
			})
			return
		}

		log.Warn("user logout")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"err": "user logout",
		})
	}
}

func AdminMiddleware(a auth.Auth, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "Auth.IsAdmin"

		log := log.With(
			slog.String("op", op),
		)

		log.Info("attempting to check users permissions")

		isAdmin := c.GetBool("IsAdmin")

		if !isAdmin {
			log.Warn("user is not admin")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"err": "user is not admin",
			})
		}

		c.Next()
	}
}
