package payments

import "gorm.io/gorm"

type PaymentsHandler struct {
	DB *gorm.DB
}

type CreatePayment struct {
	CreditAmount int    `json:"credit_amount"  validate:"required,credit_amount"`
	PriceCents   int    `json:"price_cents"  validate:"required,price_cents"`
	Provider     string `json:"provider" validate:"required,provider"`
}
