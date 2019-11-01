package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-base/middleware"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {

	router := gin.New()
	router.Use(gin.Logger())
	// router.LoadHTMLGlob("web/templates/*.tmpl.html")
	// router.Static("/web/static", "static")

	authMiddleware, err := middleware.GetJWTMiddleware()

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	middleware.SetRoutes(router, authMiddleware)

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	router.Run(":" + port)

	// router.GET("/", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{"message": "hey", "status": http.StatusOK})
	// 	//c.HTML(http.StatusOK, "index.tmpl.html", nil)
	// })

	// router.GET("/ciao", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{"message": "hey", "status": http.StatusOK})
	// })

}
