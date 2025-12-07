package internal

import (
	middlewares "github.com/GuiFernandess7/risa/internal/middlewares"
	auth "github.com/GuiFernandess7/risa/internal/modules/auth"
	filetools "github.com/GuiFernandess7/risa/internal/modules/filetools"
	payments "github.com/GuiFernandess7/risa/internal/modules/payments"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitRoutes(db *gorm.DB, e *echo.Echo) {
	authHandlers := &auth.AuthHandler{DB: db}
	paymentHandlers := &payments.PaymentsHandler{DB: db}
	e.POST("/signup", authHandlers.SignupHandler)
	e.POST("/login", authHandlers.LoginHandler)
	e.POST("/refresh", authHandlers.RefreshHandler)
	e.POST("/v1/payments/webhook/:provider", paymentHandlers.WebhookHandler)

	fileHandlers := &filetools.ImageHandler{DB: db}
	v1 := e.Group("/v1")
	v1.Use(middlewares.AuthMiddleware())
	v1.POST("/payments/create", paymentHandlers.CreatePayment)
	v1.GET("/payments/:order_id/status", paymentHandlers.GetPaymentStatus)
	v1.GET("/payments/history", paymentHandlers.GetPaymentHistory)
	v1.POST("/image/upload", fileHandlers.UploadImage)
	v1.GET("/image/status", fileHandlers.CheckStatusAsync)
}
