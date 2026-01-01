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
	"gorm.io/datatypes"
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
	fmt.Println("[WEBHOOK] - Handler: Processing request...")
	const MaxBodyBytes = int64(65536)
	req := c.Request()
	res := c.Response()

	req.Body = http.MaxBytesReader(res, req.Body, MaxBodyBytes)
	payload, err := io.ReadAll(req.Body)
	fmt.Println("[WEBHOOK] - Handler: Reading body information...")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		return c.NoContent(http.StatusServiceUnavailable)
	}

	fmt.Println("[WEBHOOK] - Handler: Getting Signature...")
	signature := req.Header.Get("Stripe-Signature")
	event, err := stripe.GetPaymentEvent(payload, signature, os.Getenv("STRIPE_WEBHOOK_SECRET"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR getting signature: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid signature",
		})
	}

	fmt.Println("[WEBHOOK] - Handler: Dispatching event...")
	paymentStatus, err := stripe.DispatchStripeEvent(event)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Webhook error: %v\n", err)
		return c.NoContent(http.StatusBadRequest)
	}

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if !paymentStatus.Success {
		fmt.Fprintf(os.Stderr, "[WEBHOOK] - Handler: Dispatch stripe error: %v\n", paymentStatus.Data)
		return c.NoContent(http.StatusOK)
	}

	fmt.Printf("[WEBHOOK] - Handler: Payment status: %v", paymentStatus)
	if paymentStatus.Type == "checkout.session.completed" {
		if paymentStatus.ProviderPaymentID == "" {
			fmt.Println("[WEBHOOK] - Handler: Provider Payment ID not found.")
			return c.NoContent(http.StatusBadRequest)
		}

		fmt.Println("[WEBHOOK] - Checkout completed successfull!")
		err := ph.handleCheckoutSessionCompleted(
			paymentStatus.ProviderPaymentID,
			payload,
		)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "error handling checkout session",
			})
		}
	}

	fmt.Println("[WEBHOOK] - Handler finished.")
	return c.NoContent(http.StatusOK)
}

func (ph PaymentsHandler) handleCheckoutSessionCompleted(
	providerPaymentID string,
	rawResponse []byte,
) error {

	tx := ph.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	paymentCrud := database.CrudGeneric[Payments]{DB: tx}
	orderCrud := database.CrudGeneric[Orders]{DB: tx}
	creditCrud := database.CrudGeneric[CreditTransactions]{DB: tx}

	payment, err := paymentCrud.FindBy("provider_payment_id", providerPaymentID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("payment not found")
	}

	if payment.Status != "pending" {
		tx.Rollback()
		return nil
	}

	payment.Status = "confirmed"
	payment.RawResponse = datatypes.JSON(rawResponse)

	if err := paymentCrud.Update(payment.ID, payment); err != nil {
		tx.Rollback()
		return err
	}

	order, err := orderCrud.FindBy("id", payment.OrderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	order.Status = "paid"
	if err := orderCrud.Update(order.ID, order); err != nil {
		tx.Rollback()
		return err
	}

	creditTx := CreditTransactions{
		UserID:      uint(order.UserID),
		Amount:      order.CreditAmount,
		Type:        "purchase",
		ReferenceID: uint(payment.ID),
		Description: "Purchase confirmed by Stripe",
	}

	if err := creditCrud.Create(&creditTx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
