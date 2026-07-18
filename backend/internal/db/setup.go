package db

import (

	"os"

	"gorm.io/gorm"
)

type Savings struct {
	id string
	amount float64
}

func Setup(db gorm.DB) {

	db.Migrator().DropTable()
	db.AutoMigrate(&Savings{})

	db.Create(&Savings{os.Getenv("AMOUNT_ID"), 0.00})
}