package main

import (
	"fmt"
	"log"
	"medical-pager-server/handlers/auth_handler"
	"net/http"
	"os"
	"time"

	stream_chat "github.com/GetStream/stream-chat-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var streamClient *stream_chat.Client

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")
	streamClient, err = stream_chat.NewClient(apiKey, []byte(apiSecret))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"access-control-allow-origin, access-control-allow-headers, Origin, X-Requested-With, Content-Type, Accept"}
	config.AllowMethods = []string{"POST", "OPTIONS"}

	router := gin.Default()
	router.POST("/auth/sign-up", auth_handler.SignUp)
	router.POST("/auth/sign-in", SignIn)

	router.Use(cors.New(config))
	router.Run()
}

func SignIn(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")

	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := c.ShouldBindJSON(&request)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}

	if request.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "username empty",
		})
		return
	}

	users, err := streamClient.QueryUsers(&stream_chat.QueryOption{
		Filter: map[string]interface{}{
			"id": request.Username,
		},
	})
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
	}

	token, err := streamClient.CreateToken(request.Username, time.Now().UTC().Add(time.Hour))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": string(token),
		"id":    request.Username,
		"name":  users[0].Name,
	})
}
