package server

import (
	"fmt"
	"goe2e-example/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RunServer() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/persons", personPostRequestHandler)
	secure := r.Group("/secure")
	secure.Use(authRequestHandler)
	{
		secure.POST("/persons", personPostRequestHandler)
	}
	r.Run()
}

func authRequestHandler(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	if auth != "john" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
}

func personPostRequestHandler(c *gin.Context) {
	var p model.Person
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusAccepted, p)
}
