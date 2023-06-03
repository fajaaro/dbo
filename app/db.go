package app

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const DB_USERNAME = "postgres"
const DB_PASSWORD = "admin"
const DB_NAME = "dbo"
const DB_HOST = "127.0.0.1"
const DB_PORT = "5432"

var db *gorm.DB

func InitDb() *gorm.DB {
	db = connectDB()
	return db
}

func GetDb() *gorm.DB {
	return db
}

func connectDB() *gorm.DB {
	conn := "host=" + DB_HOST + " user=" + DB_USERNAME + " password=" + DB_PASSWORD + " dbname=" + DB_NAME + " port=" + DB_PORT + " sslmode=disable TimeZone=Asia/Jakarta"
	fmt.Println("conn : ", conn)
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})

	if err != nil {
		fmt.Println("Error connecting to database : error=" + err.Error())
		return nil
	}

	return db
}
