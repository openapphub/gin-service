package api

import (
	"openapphub/internal/service"
	"openapphub/pkg/serializer"

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
	s := sessions.Default(c)
	s.Clear()
	s.Save()
	c.JSON(200, serializer.Response{
		Code: 0,
		Msg:  "登出成功",
	})
}
