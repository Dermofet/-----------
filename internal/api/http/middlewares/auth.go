package middlewares

import (
	"fmt"
	"music-backend-test/cmd/music-backend-test/config"
	"music-backend-test/internal/entity"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg, err := config.GetAppConfig()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if cfg.ApiKey == "" {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if len(tokenString) == 0 {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid token format"))
		}

		id, err := entity.ParseToken(tokenString)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid token: %v", err))
			return
		}

		c.Set("user-id", id)
		c.Next()
	}
}
