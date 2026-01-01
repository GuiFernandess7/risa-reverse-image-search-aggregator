package middlewares

import (
	auth "github.com/GuiFernandess7/risa/internal/modules/auth"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var jwtSecret = []byte("JWT_SECRET_KEY")

func AuthMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: jwtSecret,
		ContextKey: "user",
	})
}

func LoadUserMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			token, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return echo.ErrUnauthorized
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return echo.ErrUnauthorized
			}

			uid, ok := claims["user_id"]
			if !ok {
				return echo.ErrUnauthorized
			}

			userID := uint(uid.(float64))
			var user auth.User
			if err := db.First(&user, userID).Error; err != nil {
				return echo.ErrUnauthorized
			}

			c.Set("user", &user)
			return next(c)
		}
	}
}
