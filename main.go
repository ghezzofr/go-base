package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-base/database"
	"github.com/heroku/go-base/middleware"
	"github.com/heroku/go-base/models"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {

	database.Init()
	database.MigrateTables(&models.User{})
	defer database.Close()

	go createAdminUser()

	router := gin.New()
	router.Use(gin.Logger())
	// router.LoadHTMLGlob("web/templates/*.tmpl.html")
	// router.Static("/web/static", "static")

	authMiddleware, err := middleware.GetJWTMiddleware()

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	middleware.SetRoutes(router, authMiddleware, middleware.GetAdminHandler())

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	router.Run(":" + port)
}

func createAdminUser() {
	admin := models.User{
		Name:     "admin",
		Surname:  "admin",
		Email:    "admin@company.com",
		Password: "fantasticPassword",
	}
	admin.SetAdmin()
	_, err := admin.Create()
	log.Println(err)
}
