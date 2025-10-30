package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const ctxUserID = "userID"
const ctxUserRole = "userRole"

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "missing bearer token"})
			return
		}
		tokenStr := strings.TrimSpace(h[len("Bearer "):])
		claims, err := ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "invalid token"})
			return
		}
		c.Set(ctxUserID, claims.UserID)
		c.Set(ctxUserRole, claims.Role)
		c.Next()
	}
}

// helpers if needed by handlers
func GetUserID(c *gin.Context) (uint, bool) {
	v, ok := c.Get(ctxUserID)
	if !ok {
		return 0, false
	}
	id, _ := v.(uint)
	return id, true
}
func GetUserRole(c *gin.Context) (string, bool) {
	v, ok := c.Get(ctxUserRole)
	if !ok {
		return "", false
	}
	role, _ := v.(string)
	return role, true
}
