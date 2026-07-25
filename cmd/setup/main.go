package main

import (
	
	"log"
	"os"
    appdb "backend/internal/db" 
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
func main() {
    godotenv.Load()
    dsn := os.Getenv("DATABASE_URL")
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    appdb.Setup(*db)

}