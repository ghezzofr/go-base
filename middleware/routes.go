package middleware

import (
	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/heroku/go-base/models"
)

// SetRoutes set the entry points of the web application
func SetRoutes(router *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) {

	// ###### V1 ########
	v1 := router.Group("/v1")
	{
		// ###### AUTH #########
		v1.POST("/login", authMiddleware.LoginHandler)
		v1.GET("hi", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "hey", "status": http.StatusOK})
		})
		// v1.POST("/signin", controllers.SignIn)

		// ####### AUTH REQUIRED API #########
		auth := v1.Group("/auth")
		// Refresh time can be longer than token timeout
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
		auth.Use(authMiddleware.MiddlewareFunc())
		{
			auth.GET("/hello", func(c *gin.Context) {
				claims := jwt.ExtractClaims(c)
				user, _ := c.Get(identityKey)
				c.JSON(200, gin.H{
					"userID":   claims[identityKey],
					"userName": user.(models.User).Name,
					"text":     "Hello World.",
				})
			})
		}

	}

	router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
}
