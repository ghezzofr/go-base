package models

import (
	"log"

	"github.com/heroku/go-base/database"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User represent a User into the database
type User struct {
	gorm.Model
	Name     string `form:"name" json:"name" binding:"required"`
	Surname  string `form:"surname" json:"surname" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required" gorm:"not null; unique_index;"`
	Password string `form:"password" json:"password" binding:"required" gorm:"not null;"`
}

// IsValid makes integrity check on the current User and return a boolean with the check result
func (c *User) IsValid() (bool, error) {
	// TODO: Check uniqueness
	return true, nil
}

func (c *User) encryptPassword() {
	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	c.Password = string(hash)
}

// Create creates a User and put it into the database, if an error occours during this process an error is returned
func (c *User) Create() (bool, error) {
	// Encript the password
	c.encryptPassword()
	if err := database.DB.Create(&c); err.Error != nil && database.IsUniqueConstraintError(err.Error) {
		return false, err.Error
	}
	database.DB.Save(&c)
	return true, nil
}

// GetUserByID return the User that matches with the passed email, return nil if no User has that email
func GetUserByID(id uint) (c User, err error) {
	if dbResponse := database.DB.Where("ID = ?", id).First(&c); dbResponse.Error != nil {
		err = dbResponse.Error
	}
	return
}

// GetUserByEmail return the User that matches with the passed email, return nil if no User has that email
func GetUserByEmail(email string) (c User, err error) {
	if dbResponse := database.DB.Where("email = ?", email).First(&c); dbResponse.Error != nil {
		err = dbResponse.Error
	}
	return
}

// Authenticate return a User if the authentication process works fine and the mail and password passed matches with a DB record
func Authenticate(email string, password string) (c User, e error) {
	if c, e = GetUserByEmail(email); e != nil {
		return
	}
	_, e = checkPassword(password, c.Password)
	return
}

// checkPassword check a DB hash with a clear password and return a boolean if the password matches or not
func checkPassword(UserPassword string, hash string) (bool, error) {
	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(UserPassword)); err != nil {
		// TODO: Properly handle error
		return false, err
	}
	return true, nil
}
