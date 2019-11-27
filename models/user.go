package models

import (
	"errors"
	"log"

	"github.com/heroku/go-base/database"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Grant rapresent the kind of grant type a User can have
type Grant int

const (
	// Base is the basic user
	Base Grant = iota
	// Admin is the admin user
	Admin
)

// BasicUser represent a user without password
type BasicUser struct {
	ID      uint
	Name    string `form:"name" json:"name" binding:"required"`
	Surname string `form:"surname" json:"surname" binding:"required"`
	Email   string `form:"email" json:"email" binding:"required" gorm:"not null; unique_index;"`
	Grant   Grant
}

// User represent a User into the database
type User struct {
	gorm.Model
	Name     string `form:"name" json:"name" binding:"required"`
	Surname  string `form:"surname" json:"surname" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required" gorm:"not null; unique_index;"`
	Password string `form:"password" json:"password" binding:"required" gorm:"not null;"`
	Grant    Grant  `gorm:"not null"`
}

// IsValid makes integrity check on the current User and return a boolean with the check result
func (u *User) IsValid() (bool, error) {
	var user User
	var i int
	database.DB.Where("email = ?", u.Email).First(&user).Count(&i)
	if i == 1 {
		return false, errors.New("Duplicate value for the passed email")
	}
	return true, nil
}

func (u *User) encryptPassword() {
	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	u.Password = string(hash)
}

// Create creates a User and put it into the database, if an error occours during this process an error is returned
func (u *User) Create() (bool, error) {
	if valid, err := u.IsValid(); valid == false {
		return valid, err
	}
	// Encript the password
	u.encryptPassword()
	if err := database.DB.Create(&u); err.Error != nil && database.IsUniqueConstraintError(err.Error) {
		return false, err.Error
	}
	database.DB.Save(&u)
	return true, nil
}

// SetAdmin set the user as an administrator
func (u *User) SetAdmin() {
	u.Grant = Admin
}

// SetAdmin set the user as an administrator
func (u *User) GetBasicUser() (bu BasicUser) {
	bu.Name = u.Name
	bu.Surname = u.Surname
	bu.ID = u.ID
	bu.Grant = u.Grant
	bu.Email = u.Email
	return
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
