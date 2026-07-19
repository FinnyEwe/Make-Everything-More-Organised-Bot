package store

import (
	"backend/internal/model"
	"os"
	"gorm.io/gorm"
)

type Store struct {
	Db *gorm.DB
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) GetSavings() model.Savings {
	var savings model.Savings
	key := os.Getenv("AMOUNT_ID")
	s.Db.First(&savings, key)
	return savings
}

type Operand string
// Declare constants of this new type
const (
	OperandAdd Operand = "add"
	OperandSub Operand = "sub"
)
func (s *Store) UpdateSavings(operand Operand, amount float64) (model.Savings, error) {
	var savings model.Savings
	key := os.Getenv("AMOUNT_ID")

	operatorSymb := ""
	if operand == OperandAdd {
		operatorSymb = "+"
	} else if operand == OperandSub {
		operatorSymb = "-"
	}
	
	result := s.Db.Model(&model.Savings{}).Where("Id = "+key).Update("Amount", gorm.Expr("Amount "+ operatorSymb + " ?", amount))
	if result.Error != nil {
		return savings, result.Error
	}
	
	s.Db.First(&savings, key)
	return  savings, nil
}