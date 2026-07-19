package store

import (
	"fmt"
	"os"

	"backend/internal/model"

	"gorm.io/gorm"
)

type Store struct {
	Db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{Db: db}
}

func (s *Store) GetSavings() model.Savings {
	var savings model.Savings
	key := os.Getenv("AMOUNT_ID")
	s.Db.First(&savings, "id = ?", key)
	return savings
}

type Operand string

const (
	OperandAdd Operand = "add"
	OperandSub Operand = "sub"
)

func (s *Store) UpdateSavings(operand Operand, amount float64) (model.Savings, error) {
	var savings model.Savings
	key := os.Getenv("AMOUNT_ID")

	var operatorSymb string
	switch operand {
	case OperandAdd:
		operatorSymb = "+"
	case OperandSub:
		operatorSymb = "-"
	default:
		return savings, fmt.Errorf("unknown operand: %s", operand)
	}

	result := s.Db.Model(&model.Savings{}).
		Where("id = ?", key).
		Update("amount", gorm.Expr("amount "+operatorSymb+" ?", amount))
	if result.Error != nil {
		return savings, result.Error
	}

	s.Db.First(&savings, "id = ?", key)
	return savings, nil
}
