package middlewares

import (
	"context"
	"fmt"
	"music-backend-test/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// The NewCheckAdminMiddleware function is a middleware that checks if the user is an admin.
func NewCheckRoleMiddleware(roles []string, userInteractor usecase.UserInteractor) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, exists := c.Get("user-id")
		if !exists {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx := context.Background()
		user, err := userInteractor.GetById(ctx, userId.(uuid.UUID))
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("can't get user: %w", err))
			return
		}

		if !contains(roles, user.Role) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

func contains(s []string, v string) bool {
	for _, s_ := range s {
		if v == s_ {
			return true
		}
	}
	return false
}
