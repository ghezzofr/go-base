package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/heroku/go-base/models"
)

// SignIn create a user if not exists
func SignIn(c *gin.Context) {
	var customer models.User
	// Bind Customer by URLencoded or by JSON
	if c.Bind(&customer) != nil {
		c.BindJSON(&customer)
	}
	if isValid, err := customer.IsValid(); isValid {
		if _, err := customer.Create(); err != nil {
			c.JSON(200, err)
		} else {
			c.JSON(200, &customer)
		}
	} else {
		c.JSON(200, err)
	}
}
