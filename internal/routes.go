package internal

import (
	middlewares "github.com/GuiFernandess7/risa/internal/middlewares"
	auth "github.com/GuiFernandess7/risa/internal/modules/auth"
	filetools "github.com/GuiFernandess7/risa/internal/modules/filetools"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitRoutes(db *gorm.DB, e *echo.Echo) {
	e.POST("/signup", auth.SignupHandler)
	e.POST("/login", auth.LoginHandler)
	e.POST("/refresh", auth.RefreshHandler)

	handlers := &filetools.ImageHandler{DB: db}
	v1 := e.Group("/v1")
	v1.Use(middlewares.AuthMiddleware())

	v1.POST("/image/upload", handlers.UploadImage)
	v1.GET("/image/status", handlers.CheckStatusAsync)
}
