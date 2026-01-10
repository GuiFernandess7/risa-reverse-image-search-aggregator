package auth

import (
	"errors"

	payments "github.com/GuiFernandess7/risa/internal/modules/payments"
	database "github.com/GuiFernandess7/risa/internal/repository/database"
	"gorm.io/gorm"
)

var (
	ErrCreditBalanceNotFound = errors.New("credit balance not found")
	ErrInsufficientCredits   = errors.New("insufficient credits")
)

func VerifyUserCredits(db *gorm.DB, userID uint, cost int) error {
	crud := database.CrudGeneric[payments.CreditBalance]{DB: db}

	balance, err := crud.FindBy("user_id", userID)
	if err != nil {
		return ErrCreditBalanceNotFound
	}

	if balance.Balance < uint(cost) {
		return ErrInsufficientCredits
	}

	return nil
}
