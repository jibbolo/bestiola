package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var db gorm.DB

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "./gorm.db")
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&User{}, &Match{})
}

func main() {
	router := gin.Default()
	db.LogMode(true)
	attachUserAPI(router)
	attachMatchAPI(router)

	router.StaticFile("/", "./assets/index.html")
	router.Static("/assets", "./assets")

	// Listen and server on 0.0.0.0:8080
	router.Run(":8080")

}
