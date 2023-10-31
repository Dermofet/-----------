package middlewares

import (
	"context"
	"fmt"
	"music-backend-test/internal/entity"
	"music-backend-test/internal/usecase"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

// The NewCheckAdminMiddleware function is a middleware that checks if the user is an admin.
func NewCheckRoleMiddleware(roles []string, userInteractor usecase.UserInteractor) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the user is an admin
		userId, exists := c.Get("user-id")
		if !exists {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx := context.Background()
		user, err := userInteractor.GetById(ctx, userId.(*entity.UserID))
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("can't get user: %w", err))
			return
		}

		if !slices.Contains(roles, user.Role) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
