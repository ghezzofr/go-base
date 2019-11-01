package database

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
)

// DB represent the DB connection througt GORM
var DB *gorm.DB
var err error

// Init connect the app with the database and put into db variable an instance of *gorm.DB. initialize the database
// adapring tables to the model
func Init() {

	DB, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Println(err)
	}

	// Set as default schema, go-base
	DB.Exec("SET search_path TO 'go-base'")

}

// Close close the connection to DB. Use int with defer when you call Init or other function that use a DB connection
func Close() {
	DB.Close()
}
