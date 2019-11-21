package database

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
)

// DB represent the DB connection througt GORM
var DB *gorm.DB
var err error

// Init connect the app with the database and put into db variable an instance of *gorm.DB. initialize the database
// adapring tables to the model
func Init() {

	log.Printf("Trying to connect to the DB %s \n", os.Getenv("DATABASE_URL"))
	DB, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Println(err)
	}

	// Set as default schema, go-base
	SetSchema("go-base")
}

// Close close the connection to DB. Use int with defer when you call Init or other function that use a DB connection
func Close() {
	DB.Close()
}

// MigrateTables to the DB
func MigrateTables(tableModel interface{}) {
	DB.AutoMigrate(tableModel)
}

// SetSchema Create the schema if not exist and set it as the default schema. Every operation made on the DB connection pass by the setted schema
// For example, if you make a `SELECT * FROM tableName` query, this would be tranlated in `SELECT * FROM schemaName.tableName`
func SetSchema(schemaName string) {
	DB.Exec("CREATE SCHEMA IF NOT EXISTS \"" + schemaName + "\";")
	DB.Exec("SET search_path TO '" + schemaName + "'")
}
