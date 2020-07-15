package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func DBConnection() map[string]interface{} {

	Connection := make(map[string]interface{})
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mysqlConn := Mysql{
		Username:     os.Getenv("DB_USERNAME"),
		Password:     os.Getenv("DB_PASSWORD"),
		Host:         os.Getenv("DB_HOST"),
		Port:         os.Getenv("DB_PORT"),
		DatabaseName: os.Getenv("DB_NAME"),
	}

	mysqlConnection := mysqlConn

	Connection["mysql"] = mysqlConnection
	return Connection
}
