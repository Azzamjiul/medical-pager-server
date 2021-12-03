package auth_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func SignIn(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"password": password,
	})
}
