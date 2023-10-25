package middlewares

import (
	"fmt"
	"go-test-grpc-http/cmd/go-test-grpc-http/config"
	"go-test-grpc-http/internal/entity"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// The NewAuthMiddleware function is a middleware that handles authentication by checking for a valid
// API key and JWT token in the Authorization header.
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
