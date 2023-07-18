package handlers

import (
	"github.com/gin-gonic/gin"
	"operation-platform/services"
)

var usersServices *services.UsersService

func init() {
	usersServices = services.NewUsersService()
}

func GetUserInfo(c *gin.Context) {
	usersServices.GetUserInfo(c)
}

func GetAllUsers(c *gin.Context) {
	usersServices.GetAllUsers(c)
}

func GetUserLogin(c *gin.Context) {
	usersServices.GetUserLogin(c)
}

func GetUserLoginCallback(c *gin.Context) {
	usersServices.GetUserLoginCallback(c)
}

func PostUserLogout(c *gin.Context) {
	usersServices.PostUserLogout(c)
}
