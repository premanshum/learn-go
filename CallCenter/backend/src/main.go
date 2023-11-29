package main

import (
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

func IndexHandler(c *gin.Context) {

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
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    userlist,
	}

	c.JSON(200, gin.H{"response": response})

}

func SetupServer() *gin.Engine {

	r := gin.Default()

	r.GET("/", IndexHandler)

	return r

}

func main() {
	SetupServer().Run()
}
