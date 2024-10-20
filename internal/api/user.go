package api

import (
	"openapphub/internal/auth"
	"openapphub/internal/model"
	"openapphub/internal/service"
	"openapphub/pkg/serializer"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// UserRegister godoc
// @Summary Register a new user
// @Description Register a new user with the provided information
// @Tags user
// @Accept json
// @Produce json
// @Param user body service.UserRegisterService true "User Registration Info"
// @Success 200 {object} serializer.Response "User registered successfully"
// @Failure 400 {object} serializer.Response "Bad request"
// @Router /user/register [post]
func UserRegister(c *gin.Context) {
	var service service.UserRegisterService
	if err := c.ShouldBind(&service); err == nil {
		res := service.Register()
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// UserLogin godoc
// @Summary Log in a user
// @Description Authenticate a user and return a token
// @Tags user
// @Accept json
// @Produce json
// @Param user body service.UserLoginService true "User Login Info"
// @Success 200 {object} serializer.Response "User logged in successfully"
// @Failure 400 {object} serializer.Response "Bad request"
// @Failure 401 {object} serializer.Response "Unauthorized"
// @Router /user/login [post]
func UserLogin(c *gin.Context) {
	var service service.UserLoginService
	if err := c.ShouldBind(&service); err == nil {
		res := service.Login(c)
		c.JSON(200, res)
	} else {
		c.JSON(200, ErrorResponse(err))
	}
}

// UserMe godoc
// @Summary Get current user information
// @Description Get information about the currently logged-in user
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} serializer.Response "User information retrieved successfully"
// @Failure 401 {object} serializer.Response "Unauthorized"
// @Router /user/me [get]
func UserMe(c *gin.Context) {
	user := CurrentUser(c)
	res := serializer.BuildUserResponse(*user)
	c.JSON(200, res)
}

// UserLogout godoc
// @Summary Log out a user
// @Description Log out the currently authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} serializer.Response "User logged out successfully"
// @Failure 401 {object} serializer.Response "Unauthorized"
// @Router /user/logout [delete]
func UserLogout(c *gin.Context) {
	authMode := os.Getenv("AUTH_MODE")
	user := CurrentUser(c)

	if user == nil {
		c.JSON(401, serializer.Response{
			Code: 401,
			Msg:  "未登录",
		})
		return
	}

	if authMode == "jwt" {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "未提供令牌",
			})
			return
		}
		err := model.DeleteJWTToken(tokenString)
		if err != nil {
			c.JSON(500, serializer.DBErr("注销失败", err))
			return
		}
	} else {
		s := sessions.Default(c)
		s.Clear()
		s.Save()
	}

	c.JSON(200, serializer.Response{
		Code: 0,
		Msg:  "登出成功",
	})
}

// Add these functions to the file

func UserLogoutAll(c *gin.Context) {
	user := CurrentUser(c)
	authMode := os.Getenv("AUTH_MODE")

	if authMode == "jwt" {
		err := model.DeleteAllJWTTokensForUser(user.ID)
		if err != nil {
			c.JSON(500, serializer.DBErr("注销所有设备失败", err))
			return
		}
	} else {
		err := model.DeleteAllSessionsForUser(user.ID)
		if err != nil {
			c.JSON(500, serializer.DBErr("注销所有设备失败", err))
			return
		}
	}

	c.JSON(200, serializer.Response{
		Code: 0,
		Msg:  "已注销所有设备",
	})
}

func UserLogoutDevice(c *gin.Context) {
	user := CurrentUser(c)
	deviceID := c.Param("device_id")
	authMode := os.Getenv("AUTH_MODE")

	if user == nil {
		c.JSON(401, serializer.Response{
			Code: 401,
			Msg:  "未登录",
		})
		return
	}

	if authMode == "jwt" {
		err := model.DeleteJWTToken(deviceID)
		if err != nil {
			c.JSON(500, serializer.DBErr("注销设备失败", err))
			return
		}
	} else {
		err := model.DeleteSession(deviceID)
		if err != nil {
			c.JSON(500, serializer.DBErr("注销设备失败", err))
			return
		}
	}

	c.JSON(200, serializer.Response{
		Code: 0,
		Msg:  "已注销指定设备",
	})
}

func UserDevices(c *gin.Context) {
	user := CurrentUser(c)
	authMode := os.Getenv("AUTH_MODE")

	var devices interface{}
	var err error

	if authMode == "jwt" {
		devices, err = model.GetActiveJWTTokensForUser(user.ID)
	} else {
		devices, err = model.GetActiveSessionsForUser(user.ID)
	}

	if err != nil {
		c.JSON(500, serializer.DBErr("获取设备列表失败", err))
		return
	}

	c.JSON(200, serializer.Response{
		Code: 0,
		Data: devices,
	})
}

// Add this function to the file

// RefreshToken godoc
// @Summary Refresh JWT token
// @Description Refresh the JWT access token using a refresh token
// @Tags user
// @Accept json
// @Produce json
// @Param refresh_token body string true "Refresh Token"
// @Success 200 {object} serializer.Response "New access token"
// @Failure 400 {object} serializer.Response "Bad request"
// @Failure 401 {object} serializer.Response "Unauthorized"
// @Router /user/refresh [post]
func RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, ErrorResponse(err))
		return
	}

	newAccessToken, err := auth.RefreshToken(input.RefreshToken)
	if err != nil {
		c.JSON(401, serializer.Response{
			Code: 401,
			Msg:  "Invalid refresh token",
		})
		return
	}

	c.JSON(200, serializer.Response{
		Code: 0,
		Data: gin.H{"access_token": newAccessToken},
		Msg:  "Token refreshed successfully",
	})
}
