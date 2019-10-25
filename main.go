package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "hey", "status": http.StatusOK})
		//c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/ciao", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "hey", "status": http.StatusOK})
	})

	router.Run(":" + port)
}
