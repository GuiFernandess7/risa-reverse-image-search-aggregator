package payments

import (
	"time"

	"gorm.io/gorm"
)

type PaymentsHandler struct {
	DB *gorm.DB
}

type CreatePaymentRequest struct {
	CreditAmount int    `json:"credit_amount"  validate:"required,credit_amount"`
	Provider     string `json:"provider" validate:"required,provider"`
}

type PaymentHistoryResponse struct {
	ID           int       `json:"id"`
	CreditAmount int       `json:"credit_amount"  validate:"required,credit_amount"`
	PriceCents   int       `json:"price_cents"  validate:"required,price_cents"`
	Status       string    `json:"status" validate:"required,status"`
	CreatedAt    time.Time `json:"created_at"`
}

type WebhookPayload struct {
	EventType string      `json:"event_type"`
	Data      WebhookData `json:"data"`
}

type WebhookData struct {
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
}

type SimplePaymentIntent struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Customer string `json:"customer"`
}
