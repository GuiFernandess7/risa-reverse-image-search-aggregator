package auth

import (
	"errors"
	"time"

	database "github.com/GuiFernandess7/risa/internal/repository/database"
	"gorm.io/gorm"
)

type creditBalance struct {
	UserID    uint `gorm:"primaryKey"`
	Balance   uint `gorm:"not null"`
	UpdatedAt time.Time
}

func (creditBalance) TableName() string {
	return "credit_balance"
}

var (
	ErrCreditBalanceNotFound = errors.New("credit balance not found")
	ErrInsufficientCredits   = errors.New("insufficient credits")
)

func VerifyUserCredits(db *gorm.DB, userID uint, cost int) error {
	crud := database.CrudGeneric[creditBalance]{DB: db}

	balance, err := crud.FindBy("user_id", userID)
	if err != nil {
		return ErrCreditBalanceNotFound
	}

	if balance.Balance < uint(cost) {
		return ErrInsufficientCredits
	}

	return nil
}
