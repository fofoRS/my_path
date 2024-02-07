package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

var userDataBaseRepository UserDbConnector
var userService UserService
var userHandler UserHandler

func initializeExternalRouter(engine *gin.Engine) gin.IRoutes {
	return engine.Group("/external").
		POST("/signup", userHandler.RegisterNewUser)
}

func setupServer() *gin.Engine {
	mongoDatabase := registerMongoDataBase()
	registerUserHandler(&mongoDatabase)

	engine := gin.Default()
	initializeExternalRouter(engine)

	engine.GET("/info", func(c *gin.Context) {
		if locaIps, err := net.DefaultResolver.LookupIPAddr(context.Background(), "127.0.0.1"); err != nil {
			c.String(http.StatusInternalServerError, "cannot find local IP")
		} else {
			c.JSON(http.StatusOK, gin.H{"ip": locaIps[0]})
		}

	})

	externalRouter := engine.Group("/external")

	{
		externalRouter.POST("/user-registration")
	}
	return engine
}

func registerMongoDataBase() mongo.Database {
	hostName := os.Getenv("db_hostname")
	username := os.Getenv("db_username")
	password := os.Getenv("db_password")
	dbName := os.Getenv("db_name")
	if hostName == "" {
		panic("db_hostname env variable is not set")
	}
	dbClient := NewMongoDbClient(hostName, username, password)
	mongoDb := ConnectToDb(dbClient, dbName)
	log.Println("connected to mongo db")
	return *mongoDb
}

func registerUserHandler(mongoDb *mongo.Database) {
	var err error
	if userDataBaseRepository, err = NewUserMongoConnector(mongoDb); err != nil {
		panic(err)
	} else {
		withDataBaseConnector := func(userService *UserService) {
			userService.userDbConnector = userDataBaseRepository
		}
		userService = NewUserService(withDataBaseConnector)
		userHandler = NewUserHander(userService)
	}
}

type ApiErrorResponse struct {
	code    int
	message string
}

func (errorResponse ApiErrorResponse) GetResponseCode() int {
	return errorResponse.code
}

func (errorResponse ApiErrorResponse) GetResponseMessage() string {
	return errorResponse.message
}

type FieldValidationApiResponseBody struct {
	ValidationType string
	Fields         []string
	Message        string
}

func (apiResponseBody *FieldValidationApiResponseBody) AddField(field string) {
	apiResponseBody.Fields = append(apiResponseBody.Fields, field)
}

func ParseBindingValidationErrorToApiErrResponse(err validator.ValidationErrors) []FieldValidationApiResponseBody {

	fieldsByErrorTypeMap := make(map[string]FieldValidationApiResponseBody)
	for _, fieldError := range err {
		validationType := fieldError.Tag()
		if reponseBody, exists := fieldsByErrorTypeMap[validationType]; exists {
			reponseBody.AddField(fieldError.Field())
			fieldsByErrorTypeMap[validationType] = reponseBody
		} else {
			fields := make([]string, 0)
			fields = append(fields, fieldError.Field())
			fieldsByErrorTypeMap[validationType] = FieldValidationApiResponseBody{
				ValidationType: validationType,
				Fields:         fields,
				Message:        getValidationMessageByValidationTag(validationType),
			}
		}
	}
	var validationErrors []FieldValidationApiResponseBody
	for _, value := range fieldsByErrorTypeMap {
		validationErrors = append(validationErrors, value)
	}
	return validationErrors
}

func getValidationMessageByValidationTag(validationTag string) string {
	switch validationTag {
	case "required":
		return "Required field must not be empty"
	default:
		return "unknown validation error"
	}
}

func ParseDomainErrorToApiErrResponse(err DomainError) ApiErrorResponse {
	var errorResponseMap ApiErrorResponse
	switch err.getReason() {
	case UsernameAlreadyTaken:
		errorResponseMap = ApiErrorResponse{
			code:    http.StatusConflict,
			message: "username is not available"}
	case EmailAlreadyRegistered:
		errorResponseMap = ApiErrorResponse{
			code:    http.StatusConflict,
			message: "email address already registered"}
	default:
		errorResponseMap = ApiErrorResponse{code: http.StatusInternalServerError}
	}
	return errorResponseMap
}

func RunServer() {
	server := setupServer()
	server.Run(":8080")
}
