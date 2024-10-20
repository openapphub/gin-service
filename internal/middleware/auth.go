package middleware

import (
	"openapphub/internal/auth"
	"openapphub/internal/model"
	"openapphub/pkg/serializer"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CurrentUser 获取登录用户
func CurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		authMode := os.Getenv("AUTH_MODE")
		if authMode == "jwt" {
			tokenString := c.GetHeader("Authorization")
			if tokenString != "" {
				claims, err := auth.ParseToken(tokenString)
				if err == nil {
					user, err := model.GetUser(claims.UserID)
					if err == nil {
						c.Set("user", &user)
					}
				}
			}
		} else {
			session := sessions.Default(c)
			uid := session.Get("user_id")
			if uid != nil {
				user, err := model.GetUser(uid)
				if err == nil {
					c.Set("user", &user)
				}
			}
		}
		c.Next()
	}
}

// AuthRequired 需要登录
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authMode := os.Getenv("AUTH_MODE")
		var user *model.User

		GetZapLogger().Info("authMode", zap.String("authMode", authMode))
		if authMode == "jwt" {
			user = authenticateJWT(c)
		} else {
			user = authenticateSession(c)
		}

		if user == nil {
			c.JSON(401, serializer.CheckLogin())
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func authenticateJWT(c *gin.Context) *model.User {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return nil
	}

	claims, err := auth.ParseToken(tokenString)
	if err != nil {
		return nil
	}

	user, err := model.GetUser(claims.UserID)
	if err != nil {
		return nil
	}

	return &user
}

func authenticateSession(c *gin.Context) *model.User {
	s := sessions.Default(c)
	userID := s.Get("user_id")
	if userID == nil {
		return nil
	}

	user, err := model.GetUser(userID)
	if err != nil {
		return nil
	}

	return &user
}
