package middleware

import (
	"petshop/pkg/response"

	"github.com/gin-gonic/gin"
)

// RequirePermission checks that the authenticated admin holds the given permission.
// Must be chained after RequireAdmin.
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permsRaw, exists := c.Get(ContextAdminPerms)
		if !exists {
			response.Forbidden(c, "Access denied")
			c.Abort()
			return
		}

		perms, ok := permsRaw.([]string)
		if !ok {
			response.Forbidden(c, "Access denied")
			c.Abort()
			return
		}

		for _, p := range perms {
			if p == permission {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "Insufficient permissions")
		c.Abort()
	}
}
