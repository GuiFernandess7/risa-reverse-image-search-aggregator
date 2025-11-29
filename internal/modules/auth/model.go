package auth

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex;not null"`
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Password  string `gorm:"not null"`
	Status    string `gorm:"not null"`
	CreatedAt time.Time
}
