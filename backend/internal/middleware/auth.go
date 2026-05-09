package middleware

import (
	"strings"

	jwtpkg "petshop/pkg/jwt"
	"petshop/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	ContextCustomerID   = "customer_id"
	ContextAdminID      = "admin_id"
	ContextAdminPerms   = "admin_permissions"
)

// RequireCustomer validates a customer JWT and injects customer_id into context.
func RequireCustomer(jwtManager *jwtpkg.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearer(c)
		if token == "" {
			response.Unauthorized(c, "Authorization token required")
			c.Abort()
			return
		}

		claims, err := jwtManager.ValidateCustomerAccess(token)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set(ContextCustomerID, claims.CustomerID)
		c.Next()
	}
}

// RequireAdmin validates an admin JWT and injects admin_id + permissions into context.
func RequireAdmin(jwtManager *jwtpkg.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearer(c)
		if token == "" {
			response.Unauthorized(c, "Authorization token required")
			c.Abort()
			return
		}

		claims, err := jwtManager.ValidateAdminAccess(token)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set(ContextAdminID, claims.AdminID)
		c.Set(ContextAdminPerms, claims.Permissions)
		c.Next()
	}
}

func extractBearer(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(header, "Bearer ")
}
