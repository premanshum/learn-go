package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"teams/shared/dbEntity"
	"teams/shared/mongodb"

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

const APPPrefix = "ccTeams"

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

	var apiSettings struct {
		FilePath  string
		FileName  string
		MongoPort string
		MongoURI  string
		Version   string
	}

	apiSettings.FileName = os.Getenv("FILENAME")
	apiSettings.FilePath = os.Getenv("FILEPATH")
	apiSettings.Version = os.Getenv("VERSION")
	apiSettings.MongoPort = os.Getenv("MONGO_SVC_PORT")
	apiSettings.MongoURI = os.Getenv("MONGO_SERVER_URI")

	response := ApiResponse{
		Status:  "200",
		Message: "Web-app is healthy...",
		Data: []interface{}{
			apiSettings,
		},
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
		Message: APPPrefix + "Unable to retrieve data.",
	}

	mongoConn, err := getMongoConnection("team")
	if err != nil {
		response.Message += "mongo connection error: " + err.Error()
		logger.Error(response.Message)
		c.JSON(400, response)
		return
	}

	var teams []dbEntity.Team
	filter := dbEntity.Team{}

	//--------------Get the existing secret entity--------------
	findErr := mongoConn.Find(filter, &teams, nil)
	if findErr != nil {
		if findErr.Error() == "mongo: no documents in result" {
			response.Message += "Teams not found in DB"
			logger.Error(response.Message)
			c.JSON(400, response)
			return
		} else {
			response.Message += "error in reading data: " + findErr.Error()
			logger.Error(response.Message)
			c.JSON(400, response)
			return
		}
	}

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
		logger.Info("debug mode set to false")
	}

	if debugMode == "true" {
		mongoDbUri = os.Getenv("MONGO_SERVER_URI")
		if len(mongoDbUri) == 0 {
			fmt.Println("MONGO_SERVER_URI is missing.")
			return mongodb.MongodbConnection{}, errors.New("MONGO_SERVER_URI is missing.")
		}
	} else {
		mongoDbUri = os.Getenv("CC_MONGO_SERVER")
		if len(mongoDbUri) == 0 {
			fmt.Println("CC_MONGO_SERVER is missing.")
			logger.Warn("CC_MONGO_SERVER is missing.")
			return mongodb.MongodbConnection{}, errors.New("CC_MONGO_SERVER is missing")
		}
	}

	fmt.Println("mongodb-URI:" + mongoDbUri)
	logger.Info("mongodb-URI:" + mongoDbUri)
	return mongodb.MongodbConnection{Uri: mongoDbUri, DbName: "callcenter", CollectionName: collectionName}, nil
}
