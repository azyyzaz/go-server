package middleware

import (
	"strings"

	"go-server/internal/errcode"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// CasbinAuth 接口级权限校验，从 context 取 username 作为 subject
func CasbinAuth(e *casbin.Enforcer, whitelist ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, prefix := range whitelist {
			if strings.HasPrefix(c.Request.URL.Path, prefix) {
				c.Next()
				return
			}
		}

		username, exists := c.Get(ContextKeyUsername)
		if !exists {
			_ = c.Error(errcode.ErrUnauthorized.AsError())
			c.Abort()
			return
		}

		path := c.Request.URL.Path
		method := c.Request.Method

		// 去掉尾部斜杠，统一格式
		path = strings.TrimRight(path, "/")

		ok, err := e.Enforce(username, path, method)
		if err != nil || !ok {
			_ = c.Error(errcode.ErrForbidden.AsError())
			c.Abort()
			return
		}
		c.Next()
	}
}
