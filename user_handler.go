package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type credentials struct {
	email    string `binding:"required"`
	userName string `binding:"required"`
	password string `binding:"required"`
}

type UserHandler struct {
	userService UserService
}

func NewUserHander(userService UserService) UserHandler {
	return UserHandler{userService: userService}
}

func (handler UserHandler) RegisterNewUser(c *gin.Context) {
	var user UserApiModel
	if err := c.ShouldBindJSON(&user); err != nil {
		if validationError, ok := err.(validator.ValidationErrors); ok {
			validationApiErrorResponse := ParseBindingValidationErrorToApiErrResponse(validationError)
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationApiErrorResponse})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		if _, err := handler.userService.RegisterNewUser(user); err != nil {
			if domainError, ok := err.(DomainError); ok {
				apiErrorResponse := ParseDomainErrorToApiErrResponse(domainError)
				c.JSON(apiErrorResponse.GetResponseCode(), gin.H{"message": apiErrorResponse.message})
			} else {
				c.Status(http.StatusInternalServerError)
			}
		} else {
			c.Redirect(http.StatusPermanentRedirect, "/login")
		}
	}
}

func (handler UserHandler) Authenticate(c *gin.Context) {
	var cred credentials
	if err := c.ShouldBindJSON(&cred); err != nil {
		if validationError, ok := err.(validator.ValidationErrors); ok {
			validationApiErrorResponse := ParseBindingValidationErrorToApiErrResponse(validationError)
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationApiErrorResponse})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	} else {
		if accessToken, err := handler.userService.Authenticate(cred.email, cred.userName, cred.password); err != nil {
			if IsAuthenticationErrorReason(err) ||  {
				c.Status(http.StatusUnauthorized)
			} else {
				log.Printf("error authenticating user %s", cred.userName)
			}
		} else {

		}
	}
}
