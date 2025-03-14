package models

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var Db *gorm.DB

func init() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)
	fmt.Println(dbURI)

	Conn, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}

	Db = Conn
	err = Conn.Debug().AutoMigrate(
		&Video{}, &Stat{}, &Shedule{},
		&Session{}, &EmailConfirm{},
		&Admin{}, &Username{}, &PublicStat{},
		&User{})
	if err != nil {
		log.Fatal(err)
	}

}
