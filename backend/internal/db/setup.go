package db

import (
	"os"
	"backend/internal/model" 
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)


func Setup(db gorm.DB) {
	godotenv.Load()
	db.Migrator().DropTable()
	db.AutoMigrate(&model.Savings{})
	db.Create(&model.Savings{Id: os.Getenv("AMOUNT_ID"), Amount: 0.00})
}