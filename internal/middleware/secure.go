// secure.go
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

// SecureMiddleware 返回一个 Gin 中间件，用于设置安全相关的 HTTP 头
func SecureMiddleware() gin.HandlerFunc {
	secureMiddleware := secure.New(secure.Options{
		// SSLRedirect:           true,
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
	})

	return func(c *gin.Context) {
		// 排除特定路由的 CSP 策略
		if c.FullPath() == "/swagger/*any" {
			c.Next()
			return
		}

		err := secureMiddleware.Process(c.Writer, c.Request)
		if err != nil {
			c.Abort()
			return
		}
		c.Next()
	}
}
