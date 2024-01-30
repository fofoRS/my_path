package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

type UserDbConnector interface {
	CreateNewUser(User) (EntityId, error)
	UsernameExist(string) (bool, error)
	UserRegistered(string) (bool, error)
}

type UserMongoDbConnector struct {
	collection mongo.Collection
}

func NewUserMongoConnector(clientDb *mongo.Database) (UserDbConnector, error) {
	if clientDb == nil {
		return nil, errors.New("clientDb Argument cannot be nil")
	}

	return &UserMongoDbConnector{
		collection: *clientDb.Collection("users"),
	}, nil
}

func (mongoConnector UserMongoDbConnector) CreateNewUser(newUser User) (EntityId, error) {
	result, err := mongoConnector.collection.InsertOne(context.Background(), newUser)
	if err != nil {
		return "", err
	} else {
		objectId, ok := result.InsertedID.(primitive.ObjectID)
		if ok {
			return EntityId(objectId.String()), nil
		} else {
			return "", errors.New("Cannot convert id to ObjectId type")
		}
	}
}

func (mongoConnector UserMongoDbConnector) UsernameExist(username string) (bool, error) {
	result := mongoConnector.collection.FindOne(context.Background(), bson.D{{Key: "username", Value: username}})
	if result.Err() != nil && result.Err() == mongo.ErrNoDocuments {
		return false, nil
	} else if result.Err() != nil {
		return false, result.Err()
	}

	var user *User
	if err := result.Decode(user); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (mongoConnector UserMongoDbConnector) UserRegistered(email string) (bool, error) {
	result := mongoConnector.collection.FindOne(context.Background(), bson.D{{Key: "email", Value: email}})
	if result.Err() != nil && result.Err() == mongo.ErrNoDocuments {
		return false, nil
	} else if result.Err() != nil {
		return false, result.Err()
	}

	var user *User
	if err := result.Decode(user); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

type UserService struct {
	userDbConnector UserDbConnector
}

func NewUserService(constuctorOptions ...func(*UserService)) *UserService {
	userService := &UserService{}
	for _, option := range constuctorOptions {
		option(userService)
	}
	return userService
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
