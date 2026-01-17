package db

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() {
	log.Printf("Connecting to %s", dbHost())

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUsername(), dbPassword(), dbHost(), dbPort(), dbName())

	db0, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db = db0
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Successfully connected to database")
}

func dbHost() string {
	val, ok := os.LookupEnv("DB_HOST")
	if !ok {
		return "localhost"
	}
	return val
}

func dbPort() int {
	val, ok := os.LookupEnv("DB_PORT")
	if !ok {
		return 3306
	}
	port, err := strconv.Atoi(val)
	if err != nil {
		log.Fatal("Failed to parse DB_PORT:", err)
	}
	return port
}

func dbName() string {
	val, ok := os.LookupEnv("DB_NAME")
	if !ok {
		return "omniscience"
	}
	return val
}

func dbUsername() string {
	val, ok := os.LookupEnv("DB_USERNAME")
	if !ok {
		return "admin"
	}
	return val
}

func dbPassword() string {
	val, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		return "admin"
	}
	return val
}

func SetDB(database *gorm.DB) {
	db = database
}
