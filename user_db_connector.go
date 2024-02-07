package main

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDbConnector interface {
	CreateNewUser(User) (EntityId, error)
	UsernameExist(string) (bool, error)
	UserRegistered(string) (bool, error)
	ValidateCredentials(userName, password string) (*User, error)
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

	var user User
	if err := result.Decode(&user); err != nil {
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

func (mongoConnector UserMongoDbConnector) ValidateCredentials(username, password string) (*User, error) {
	result := mongoConnector.collection.FindOne(context.Background(), bson.D{{Key: "username", Value: username}, {Key: "password", Value: password}})
	if result.Err() != nil && result.Err() == mongo.ErrNoDocuments {
		return &User{}, nil
	} else if result.Err() != nil {
		return &User{}, result.Err()
	}
	var user *User
	if err := result.Decode(user); err != nil {
		log.Printf("Error decoding result, %s", err.Error())
		return &User{}, err
	}

	return user, nil
}
