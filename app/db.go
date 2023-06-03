package app

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDb() *gorm.DB {
	db = connectDB()
	return db
}

func GetDb() *gorm.DB {
	return db
}

func connectDB() *gorm.DB {
	var DB_USERNAME = os.Getenv("DB_USERNAME")
	var DB_PASSWORD = os.Getenv("DB_PASSWORD")
	var DB_NAME = os.Getenv("DB_NAME")
	var DB_HOST = os.Getenv("DB_HOST")
	var DB_PORT = os.Getenv("DB_PORT")

	conn := "host=" + DB_HOST + " user=" + DB_USERNAME + " password=" + DB_PASSWORD + " dbname=" + DB_NAME + " port=" + DB_PORT + " sslmode=disable TimeZone=Asia/Jakarta"
	fmt.Println("conn : ", conn)
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})

	if err != nil {
		fmt.Println("Error connecting to database : error=" + err.Error())
		return nil
	}

	return db
}
