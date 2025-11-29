package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	database "github.com/GuiFernandess7/risa/internal/repository/database"
)

var jwtSecret = []byte("JWT_SECRET_KEY")

func (ah AuthHandler) LoginHandler(c echo.Context) error {
	var body LoginRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(400, echo.Map{"error": "invalid body"})
	}

	if err := c.Validate(&body); err != nil {
		return c.JSON(400, echo.Map{"error": "invalid fields"})
	}

	crud := database.CrudGeneric[User]{DB: ah.DB}
	user, err := crud.FindBy("email", body.Email)
	if err != nil {
		return c.JSON(400, echo.Map{"error": "invalid credentials"})
	}

	if !CheckPasswordHash(body.Password, user.Password) {
		return c.JSON(400, echo.Map{"error": "invalid credentials"})
	}

	access, refresh, err := generateTokens(user.ID)
	if err != nil {
		return c.JSON(500, echo.Map{"error": "token error"})
	}

	return c.JSON(200, echo.Map{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (ah AuthHandler) SignupHandler(c echo.Context) error {
	var body SignupRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(400, echo.Map{"error": "invalid body"})
	}

	if err := c.Validate(&body); err != nil {
		return c.JSON(400, echo.Map{"error": "invalid fields"})
	}

	crud := database.CrudGeneric[User]{DB: ah.DB}
	_, err := crud.FindBy("email", body.Email)
	if err == nil {
		return c.JSON(400, echo.Map{"error": "email already exists"})
	}

	hashedPwd, err := HashPassword(body.Password)
	if err != nil {
		return c.JSON(500, echo.Map{"error": "password hash error"})
	}

	newUser := User{
		Email:     body.Email,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Password:  hashedPwd,
		Status:    "active",
	}

	if err := crud.Create(&newUser); err != nil {
		return c.JSON(500, echo.Map{"error": "database error"})
	}

	return c.JSON(201, echo.Map{
		"message": "Account created successfully",
		"user_id": newUser.ID,
	})
}

func (ah AuthHandler) RefreshHandler(c echo.Context) error {
	var body RefreshTokenRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(400, echo.Map{"error": "invalid body"})
	}

	if err := c.Validate(&body); err != nil {
		return c.JSON(400, echo.Map{"error": "invalid fields"})
	}

	token, err := jwt.Parse(body.RefreshToken, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.JSON(401, echo.Map{"error": "invalid refresh token"})
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["type"] != "refresh" {
		return c.JSON(401, echo.Map{"error": "not a refresh token"})
	}

	userID := int(claims["user_id"].(float64))

	access, refresh, err := generateTokens(uint(userID))
	if err != nil {
		return c.JSON(500, echo.Map{"error": "token generation failed"})
	}

	return c.JSON(200, echo.Map{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func generateTokens(userID uint) (accessToken string, refreshToken string, err error) {
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}

	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"type":    "refresh",
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessToken, err = access.SignedString(jwtSecret)
	if err != nil {
		return
	}

	refreshToken, err = refresh.SignedString(jwtSecret)
	return
}
