package services

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"net/http"
	"operation-platform/utils"
	"os"
)

var (
	AzureADConfig *oauth2.Config
	SessionStore  *sessions.CookieStore
)

func init() {
	AzureADConfig = &oauth2.Config{
		ClientID:     os.Getenv("AZURE_CLIENT_ID"),
		ClientSecret: os.Getenv("AZURE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AZURE_REDIRECT_URL"),
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
			TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
		},
	}
}

func NewUsersService() *UsersService {
	return &UsersService{}
}

type UsersService struct {
}

func (s *UsersService) GetUserInfo(c *gin.Context) {

	c.JSON(http.StatusOK, utils.Response{
		Code:    utils.SuccessCode,
		Message: utils.SuccessMessage,
		Data:    "",
	})

}

func (s *UsersService) getUserInfo(userId int64) (interface{}, error) {
	return userId, nil
}

func (s *UsersService) GetAllUsers(c *gin.Context) {

	c.JSON(http.StatusOK, utils.Response{
		Code:    utils.SuccessCode,
		Message: utils.SuccessMessage,
		Data:    "",
	})
}

func (s *UsersService) getAllUsers() ([]string, error) {
	return nil, nil
}

func (s *UsersService) GetUserLogin(c *gin.Context) {
	url := AzureADConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

func (s *UsersService) GetUserLoginCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := AzureADConfig.Exchange(c.Request.Context(), code)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	accessToken := token.AccessToken
	idToken := token.Extra("id_token")

	c.JSON(http.StatusOK, utils.Response{
		Code:    utils.SuccessCode,
		Message: utils.SuccessMessage,
		Data:    gin.H{"access_token": accessToken, "id_token": idToken},
	})
}

func (s *UsersService) PostUserLogout(c *gin.Context) {

	c.JSON(http.StatusOK, utils.Response{
		Code:    utils.SuccessCode,
		Message: utils.SuccessMessage,
		Data:    "",
	})
}

func (s *UsersService) postUserLogout() (interface{}, error) {
	return nil, nil
}
