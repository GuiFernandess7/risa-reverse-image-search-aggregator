package payments

import (
	"time"

	"gorm.io/datatypes"
)

type Orders struct {
	ID           uint   `gorm:"primaryKey"`
	UserID       int    `gorm:"not null"`
	CreditAmount int    `gorm:"not null"`
	PriceCents   int    `gorm:"not null"`
	Status       string `gorm:"not null"`
	CreatedAt    time.Time
}

type Payments struct {
	ID                uint   `gorm:"primaryKey"`
	OrderID           int64  `gorm:"not null"`
	Provider          string `gorm:"not null"`
	ProviderPaymentID string `gorm:"not null"`
	Status            string `gorm:"not null"`
	RawResponse       datatypes.JSON
	CreatedAt         time.Time
}

type CreditTransactions struct {
	UserID      uint   `gorm:"primaryKey"`
	Amount      int   `gorm:"not null"`
	Type        string `gorm:"not null"`
	ReferenceID uint   `gorm:"not null"`
	Description string `gorm:"not null"`
}
