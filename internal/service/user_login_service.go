package service

import (
	"openapphub/internal/auth"
	"openapphub/internal/model"
	"openapphub/internal/util"
	"openapphub/pkg/serializer"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// UserLoginService 管理用户登录的服务
type UserLoginService struct {
	UserName   string `form:"user_name" json:"user_name" binding:"required,min=5,max=30"`
	Password   string `form:"password" json:"password" binding:"required,min=8,max=40"`
	DeviceInfo string `form:"device_info" json:"device_info"`
}

func (service *UserLoginService) Login(c *gin.Context) serializer.Response {
	var user model.User

	if err := model.DB.Where("user_name = ?", service.UserName).First(&user).Error; err != nil {
		return serializer.ParamErr("账号或密码错误", nil)
	}

	if !user.CheckPassword(service.Password) {
		return serializer.ParamErr("账号或密码错误", nil)
	}

	authMode := os.Getenv("AUTH_MODE")
	if authMode == "jwt" {
		return service.loginWithJWT(c, user)
	} else {
		return service.loginWithSession(c, user)
	}
}

func (service *UserLoginService) loginWithJWT(_ *gin.Context, user model.User) serializer.Response {
	accessToken, refreshToken, err := auth.GenerateTokenPair(user)
	if err != nil {
		return serializer.Err(serializer.CodeEncryptError, "生成令牌失败", err)
	}

	expiresAt := time.Now().Add(time.Hour * 24) // Token expires in 24 hours
	err = model.CreateJWTToken(user.ID, accessToken, service.DeviceInfo, expiresAt)
	if err != nil {
		return serializer.DBErr("保存令牌失败", err)
	}

	return serializer.BuildUserResponseWithToken(user, accessToken, refreshToken)
}

func (service *UserLoginService) loginWithSession(c *gin.Context, user model.User) serializer.Response {
	s := sessions.Default(c)
	sessionID := util.RandStringRunes(32)
	s.Set("user_id", user.ID)
	s.Set("session_id", sessionID)
	err := s.Save()
	if err != nil {
		return serializer.Err(serializer.CodeEncryptError, "保存会话失败", err)
	}

	expiresAt := time.Now().Add(time.Hour * 24 * 7) // Session expires in 7 days
	err = model.CreateSession(user.ID, sessionID, service.DeviceInfo, expiresAt)
	if err != nil {
		return serializer.DBErr("保存会话失败", err)
	}

	return serializer.BuildUserResponse(user)
}
