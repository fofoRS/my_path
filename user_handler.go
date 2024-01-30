package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService UserService
}

func NewUserHander(userService UserService) UserHandler {
	return UserHandler{userService: userService}
}

func (handler UserHandler) registerNewUser(c *gin.Context) {
	var user UserApiModel
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		if _, err := handler.userService.RegisterNewUser(user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.Redirect(http.StatusPermanentRedirect, "/login")
		}
	}
}
