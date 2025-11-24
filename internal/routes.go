package internal

import (
	filetools "github.com/GuiFernandess7/risa/internal/modules/filetools"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitRoutes(db *gorm.DB, e *echo.Echo) {
	handlers := &filetools.ImageHandler{DB: db}

	v1 := e.Group("/v1")
	v1.POST("/image/upload", handlers.UploadImage)
	v1.GET("/image/status", handlers.CheckStatusAsync)
}
