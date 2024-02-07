package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WebApiErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // execute next handle in the chain first.
		if c.IsAborted() {
			handlerError := c.Err()
			domainError, ok := handlerError.(DomainError)
			if ok {
				switch domainError.getReason() {
				case UsernameAlreadyTaken:
					c.AbortWithStatusJSON(http.StatusConflict, gin.H{"err": handlerError.Error()})
				case EmailAlreadyRegistered:
					c.AbortWithStatusJSON(http.StatusConflict, gin.H{"err": handlerError.Error()})
				}
			}
		}
	}
}
