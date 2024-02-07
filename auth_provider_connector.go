package main

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthProviderConnector interface {
	GetToken(userName, password string) (string, error)
}

type DbAuthProvider struct {
	dbConnector UserDbConnector
}

func (a DbAuthProvider) GetToken(userName, password string) (string, error) {
	if dbUser, err := a.dbConnector.ValidateCredentials(userName, password); err != nil {
		log.Println("error occurred validating user")
		return "", SystemError{error: err}
	} else if dbUser != nil {
		now := time.Now()
		unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":   userName,
			"exp":   now.Add(1 * time.Hour),
			"iat":   now,
			"roles": dbUser.Roles,
		})
		if signedToken, err := unsignedToken.SignedString([]byte(jwt.SigningMethodHS256.Hash.String())); err != nil {
			log.Print("failed signing access token", err)
			return "", err
		} else {
			return signedToken, nil
		}
	} else {
		return "", UserError{
			reason:  LoginAttemptFailed,
			message: "Username or Password are invalid",
		}
	}
}
