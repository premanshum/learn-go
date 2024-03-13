package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
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
