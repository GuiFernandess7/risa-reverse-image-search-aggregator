package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var jwtSecret = []byte("JWT_SECRET_KEY")

func LoginHandler(c echo.Context) error {
	var body LoginRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(400, echo.Map{"error": "invalid body"})
	}

	if err := c.Validate(&body); err != nil {
		return c.JSON(400, echo.Map{"error": "invalid fields"})
	}

	// database operation
	userID := 123

	access, refresh, err := generateTokens(userID)
	if err != nil {
		return c.JSON(500, echo.Map{"error": "token error"})
	}

	return c.JSON(200, echo.Map{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func SignupHandler(c echo.Context) error {
	var body SignupRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(400, echo.Map{"error": "invalid body"})
	}

	if err := c.Validate(&body); err != nil {
		return c.JSON(400, echo.Map{"error": "invalid fields"})
	}

	// database operation
	userID := 123

	return c.JSON(201, echo.Map{
		"message": "signup ok",
		"user_id": userID,
	})
}

func RefreshHandler(c echo.Context) error {
	type Body struct {
		RefreshToken string `json:"refresh_token"`
	}

	var body Body
	if err := c.Bind(&body); err != nil {
		return c.JSON(400, echo.Map{"error": "invalid body"})
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

	access, refresh, err := generateTokens(userID)
	if err != nil {
		return c.JSON(500, echo.Map{"error": "token generation failed"})
	}

	return c.JSON(200, echo.Map{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func generateTokens(userID int) (accessToken string, refreshToken string, err error) {
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
