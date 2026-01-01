package services 

import (
	auth "github.com/GuiFernandess7/risa/internal/modules/auth"
	"github.com/labstack/echo/v4"
)

func GetAuthUser(c echo.Context) (*auth.User, error) {
	user, ok := c.Get("user").(*auth.User)
	if !ok {
		return nil, echo.ErrUnauthorized
	}
	return user, nil
}

