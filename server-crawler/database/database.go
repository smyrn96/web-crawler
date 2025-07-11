package database

import (
	"backend-crawler/config"
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	once sync.Once
)

func Connect() {
	once.Do(func() {
		var err error
		DB, err = gorm.Open(mysql.Open(config.DB_DSN), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		log.Println("Connected to MySQL")
	})
}
