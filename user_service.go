package main

import (
	"fmt"
	"log"
)

type Role string

type EntityId string

const (
	RootManager Role = "root_manager"
	Manager     Role = "manager"
	Member      Role = "member"
)

type Principal struct {
	UserName string
	Roles    []Role
}

type User struct {
	EntityId  string
	Username  string
	FirstName string
	Lastname  string
	Email     string
	Password  string
	Roles     []Role
	Companies []Company
}

type UserApiModel struct {
	EntityId  string
	Username  string `binding:"required"`
	FirstName string `binding:"required"`
	Lastname  string `binding:"required"`
	Email     string `binding:"required,email"`
	Password  string `binding:"required"`
}

type UserService struct {
	userDbConnector UserDbConnector
	authProvider    AuthProviderConnector
}

func NewUserService(constuctorOptions ...func(*UserService)) UserService {
	userService := &UserService{}
	for _, option := range constuctorOptions {
		option(userService)
	}
	return *userService
}

func (service UserService) RegisterNewUser(userApiModel UserApiModel) (EntityId, error) {
	user := User{
		FirstName: userApiModel.FirstName,
		Lastname:  userApiModel.Lastname,
		Username:  userApiModel.Username,
		Password:  userApiModel.Password,
		Email:     userApiModel.Email,
	}

	if exists, err := service.userDbConnector.UsernameExist(userApiModel.Username); err != nil {
		log.Printf("Error occurred talking with DB -- %v\n", err)
		return "", UserError{
			message: "cannot check if username is already taken",
		}
	} else if exists {
		log.Printf("username %s is taken", userApiModel.Username)
		return "", UserError{
			message: fmt.Sprintf("username %s already taken", userApiModel.Username),
		}
	}

	if emailRegistered, err := service.userDbConnector.UserRegistered(userApiModel.Email); err != nil {
		log.Printf("Error occurred talking with DB --%v\n", err)
		return "", UserError{
			message: "cannot check if user is already registered",
		}
	} else if emailRegistered {
		return "", UserError{
			message: fmt.Sprintf("email %s already registered", userApiModel.Email),
		}
	}

	if id, err := service.userDbConnector.CreateNewUser(user); err != nil {
		log.Printf("error inserting user to DB, with error %s", err.Error())
		return "", UserError{
			message: fmt.Sprintf("Cannot register user %s", userApiModel.Username),
		}
	} else {
		return EntityId(id), nil
	}
}

/*
1. verify the user is registered
2. call authenticatorProvider
*/
func (service UserService) Authenticate(email, userName, password string) (string, error) {
	if isRegistered, err := service.userDbConnector.UserRegistered(email); err != nil {
		log.Print("error checking if user is registed", err)
		return "", err
	} else if isRegistered {
		if accessToken, err := service.authProvider.GetToken(userName, password); err != nil {
			return "", err
		} else {
			return accessToken, nil
		}
	} else {
		return "", UserError{
			reason:  UserNotRegistered,
			message: "user is not registed.",
		}
	}
}
