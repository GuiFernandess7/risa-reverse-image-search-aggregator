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

	e.POST("api/signup", authHandlers.SignupHandler)
	e.POST("api//login", authHandlers.LoginHandler)
	e.POST("api/refresh", authHandlers.RefreshHandler)

	e.POST("api/v1/payments/webhook/:provider", paymentHandlers.WebhookHandler)

	fileHandlers := &filetools.ImageHandler{DB: db}
	v1 := e.Group("/v1")
	v1.Use(
		middlewares.AuthMiddleware(),
		middlewares.LoadUserMiddleware(db),
	)
	v1.POST("api/payments/create", paymentHandlers.CreatePayment)
	v1.GET("api/payments/:order_id/status", paymentHandlers.GetPaymentStatus)
	v1.GET("api/payments/history", paymentHandlers.GetPaymentHistory)
	v1.POST("api/image/upload", fileHandlers.UploadImage)
	v1.GET("api/image/status", fileHandlers.CheckStatusAsync)
}
