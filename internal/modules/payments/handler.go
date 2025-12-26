package payments

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	database "github.com/GuiFernandess7/risa/internal/repository/database"
	stripe "github.com/GuiFernandess7/risa/internal/services/stripe"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var allowedProviders = []string{"stripe"}

const PricePerCreditCents = 200

func (ph PaymentsHandler) CreatePayment(c echo.Context) error {
	var body CreatePaymentRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid body request",
		})
	}

	if err := c.Validate(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid fields"})
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := int(claims["user_id"].(float64))

	orderCrud := database.CrudGeneric[Orders]{DB: ph.DB}
	order := Orders{
		UserID:       userID,
		CreditAmount: body.CreditAmount,
		PriceCents:   PricePerCreditCents,
		Status:       "pending",
	}

	if err := orderCrud.Create(&order); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not create order"})
	}

	checkoutSession, err := stripe.CreateCheckoutSession(
		int64(userID),
		int64(body.CreditAmount),
		int64(PricePerCreditCents),
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "payment error"})
	}

	paymentCrud := database.CrudGeneric[Payments]{DB: ph.DB}
	payment := Payments{
		OrderID:           int64(order.ID),
		Provider:          "stripe",
		ProviderPaymentID: checkoutSession.ID,
		Status:            "pending",
	}

	if err := paymentCrud.Create(&payment); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not create payment record"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":      "payment initiated",
		"checkout_url": checkoutSession.URL,
		"order_id":     order.ID,
		"payment_id":   payment.ID,
		"provider_id":  checkoutSession.ID,
	})
}

func (ph PaymentsHandler) GetPaymentStatus(c echo.Context) error {
	orderIDStr := c.Param("order_id")

	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid order_id",
		})
	}

	crud := database.CrudGeneric[Orders]{DB: ph.DB}
	order, err := crud.FindBy("id", orderID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "order not found",
		})
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))
	if int64(order.UserID) != userID {
		return c.JSON(http.StatusForbidden, echo.Map{
			"error": "not authorized",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"order_status":   order.Status,
		"payment_status": order.Status,
		"amount":         order.CreditAmount,
	})
}

func (ph PaymentsHandler) GetPaymentHistory(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := int(claims["user_id"].(float64))

	crud := database.CrudGeneric[Orders]{DB: ph.DB}
	orders, err := crud.Read("user_id", userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "error getting history",
		})
	}

	history := make([]PaymentHistoryResponse, len(orders))
	for i, o := range orders {
		history[i] = PaymentHistoryResponse{
			ID:           int(o.ID),
			CreditAmount: o.CreditAmount,
			PriceCents:   o.PriceCents,
			Status:       o.Status,
			CreatedAt:    o.CreatedAt,
		}
	}
	return c.JSON(http.StatusOK, history)
}

func (ph PaymentsHandler) WebhookHandler(c echo.Context) error {
	const MaxBodyBytes = int64(65536)
	req := c.Request()
	res := c.Response()

	req.Body = http.MaxBytesReader(res, req.Body, MaxBodyBytes)

	payload, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		return c.NoContent(http.StatusServiceUnavailable)
	}

	signature := req.Header.Get("Stripe-Signature")

	event, err := stripe.GetPaymentEvent(payload, signature, os.Getenv("STRIPE_WEBHOOK_SECRET"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR getting signature: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid signature",
		})
	}

	_, err = stripe.DispatchStripeEvent(event)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Webhook error: %v\n", err)
		return c.NoContent(http.StatusBadRequest)
	}

	return c.NoContent(http.StatusOK)
}
