package main

import (
	"callcenter/shared/dbEntity"
	"callcenter/shared/mongodb"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Define a complex data structure
type ApiResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    []interface{} `json:"data"`
}

// Another struct for a specific type of data
type UserData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func main() {
	SetupServer().Run()
}

func SetupServer() *gin.Engine {

	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.GET("/", HealthHandler)
	r.GET("/teams", GetTeamsHandler)
	r.GET("/teamsdb", GetTeamsDBHandler)
	return r
}

func HealthHandler(c *gin.Context) {
	fmt.Println("Env-FileName:", os.Getenv("FILEPATH")+os.Getenv("FILENAME"))
	response := ApiResponse{
		Status:  os.Getenv("FILEPATH") + os.Getenv("FILENAME"),
		Message: "Web-app is healthy...",
	}

	c.JSON(200, response)
}

func GetTeamsHandler(c *gin.Context) {
	// Open jsonFile
	jsonFile, err := os.Open(os.Getenv("FILEPATH") + os.Getenv("FILENAME"))
	// if os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var result map[string]interface{}

	json.Unmarshal([]byte(byteValue), &result)

	array := result["teams"].([]interface{})

	fmt.Println(array)

	for i := range array {
		fmt.Println(array[i])
		//fmt.Printf("key=%s value=%+v\n ", k, v)
	}

	response := ApiResponse{
		Status:  os.Getenv("FILEPATH") + os.Getenv("FILENAME"),
		Message: "Data retrieved successfully",
		Data:    array,
	}

	c.JSON(200, response)
}

func GetTeamsHandler_1(c *gin.Context) {
	var userlist []interface{}

	userlist = append(userlist,
		UserData{
			Username: "prem",
			Email:    "prem@example.com",
		})

	userlist = append(userlist,
		UserData{
			Username: "priya",
			Email:    "priya@example.com",
		})

	response := ApiResponse{
		Status:  os.Getenv("NEWVAR"),
		Message: "Data retrieved successfully",
		Data:    userlist,
	}

	c.JSON(200, response)
}

func GetTeamsDBHandler(c *gin.Context) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	response := ApiResponse{
		Status:  "400",
		Message: "Unable to retrieve data.",
	}

	mongoConn, err := getMongoConnection("team")
	if err != nil {
		errorMessage := "error in getting mongo connection: " + err.Error()
		logger.Error(errorMessage)
		response.Message += errorMessage
		c.JSON(400, response)
	}

	var teams []dbEntity.Team
	filter := dbEntity.Team{}

	//--------------Get the existing secret entity--------------
	findErr := mongoConn.Find(filter, &teams, nil)
	if findErr != nil {
		if findErr.Error() == "mongo: no documents in result" {
			errorMessage := "Teams not found in DB"
			logger.Error(errorMessage)
			response.Message += errorMessage
			c.JSON(400, response)
		} else {
			errorMessage := "error in reading secret from DB: " + findErr.Error()
			logger.Error(errorMessage)
			response.Message += errorMessage
			c.JSON(400, response)
		}
	}

	// Open jsonFile

	var array []interface{}

	for _, v := range teams {
		array = append(array, v)
	}

	response = ApiResponse{
		Status:  "200",
		Message: "Data retrieved successfully",
		Data:    array,
	}

	c.JSON(200, response)
}

func getMongoConnection(collectionName string) (mongodb.MongodbConnection, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	var mongoDbUri string
	debugMode := os.Getenv("DEBUG_MODE")
	if len(debugMode) == 0 {
		debugMode = "false"
		logger.Warn("DEBUG_MODE is not specified in the DEBUG_MODE environment variable. Hence setting it to false")
	}

	if debugMode == "true" {
		mongoDbUri = os.Getenv("MONGO_SERVER_URI")
		if len(mongoDbUri) == 0 {
			fmt.Println("Environment variable MONGO_SERVER_URI is missing.")
			return mongodb.MongodbConnection{}, errors.New("mongoDbUri must be specified in the MONGO_SERVER_URI environment variable")
		}
	}

	return mongodb.MongodbConnection{Uri: mongoDbUri, DbName: "callcenter", CollectionName: collectionName}, nil
}
