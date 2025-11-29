package middlewares

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var jwtSecret = []byte("JWT_SECRET_KEY")

func AuthMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: jwtSecret,
		ContextKey: "user",
	})
}
