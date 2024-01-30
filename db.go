package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDbClient(hostName, username, password string) *mongo.Client {
	if hostName == "" {
		panic(errors.New("hostName argument is required"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	var client *mongo.Client
	var err error
	if username == "" || password == "" {
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:27017", hostName)))
	} else {
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:27017", username, password, hostName)))
	}
	if err != nil {
		panic(err)
	} else {
		return client
	}
}

func ConnectToDb(client *mongo.Client, dbName string) *mongo.Database {
	if dbName == "" {
		dbName = "local"
	}
	return client.Database(dbName)
}
