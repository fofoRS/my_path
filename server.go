package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func setupServer() *gin.Engine {

	hostName := os.Getenv("db_hostname")
	username := os.Getenv("db_username")
	password := os.Getenv("db_password")
	dbName := os.Getenv("db_name")
	if hostName == "" {
		panic("db_hostname env variable is not set")
	}
	dbClient := NewMongoDbClient(hostName, username, password)
	_ = ConnectToDb(dbClient, dbName)
	log.Println("connected to mongo db")

	engine := gin.Default()
	engine.GET("/info", func(c *gin.Context) {
		if locaIps, err := net.DefaultResolver.LookupIPAddr(context.Background(), "127.0.0.1"); err != nil {
			c.String(http.StatusInternalServerError, "cannot find local IP")
		} else {
			c.JSON(http.StatusOK, gin.H{"ip": locaIps[0]})
		}

	})
	return engine
}

func RunServer() {
	server := setupServer()
	server.Run(":8080")
}
