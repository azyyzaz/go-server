package middleware

import (
	"strings"

	"go-server/internal/errcode"
	appjwt "go-server/internal/jwt"

	"github.com/gin-gonic/gin"
)

const (
	ContextKeyUserID   = "user_id"
	ContextKeyUsername = "username"
)

// JWTAuth JWT 鉴权中间件，whitelist 中的路径前缀直接放行
func JWTAuth(manager *appjwt.Manager, blacklist *appjwt.Blacklist, whitelist ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, prefix := range whitelist {
			if strings.HasPrefix(c.Request.URL.Path, prefix) {
				c.Next()
				return
			}
		}

		tokenStr := extractToken(c)
		if tokenStr == "" {
			_ = c.Error(errcode.ErrUnauthorized.AsError())
			c.Abort()
			return
		}

		claims, err := manager.Parse(tokenStr)
		if err != nil {
			_ = c.Error(errcode.ErrUnauthorized.AsError())
			c.Abort()
			return
		}

		if claims.TokenType != appjwt.AccessToken {
			_ = c.Error(errcode.ErrUnauthorized.AsError())
			c.Abort()
			return
		}

		if blacklist != nil {
			blocked, err := blacklist.IsBlocked(c.Request.Context(), tokenStr)
			if err != nil || blocked {
				_ = c.Error(errcode.ErrUnauthorized.AsError())
				c.Abort()
				return
			}
		}

		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUsername, claims.Username)
		c.Next()
	}
}

// extractToken 从 Authorization: Bearer <token> 中提取 token
func extractToken(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimPrefix(header, "Bearer ")
	}
	return ""
}
