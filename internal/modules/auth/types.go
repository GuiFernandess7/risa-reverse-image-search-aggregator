package auth

import "gorm.io/gorm"

type AuthHandler struct {
	DB *gorm.DB
}

type SignupRequest struct {
	Email     string `json:"email"        validate:"required,email"`
	FirstName string `json:"first_name"   validate:"required"`
	LastName  string `json:"last_name"    validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

