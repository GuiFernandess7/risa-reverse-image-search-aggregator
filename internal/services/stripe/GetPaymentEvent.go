package stripe

import (
	"encoding/json"
	"fmt"
	"os"

	stripe "github.com/stripe/stripe-go/v80"
)

type PaymentStatus struct {
	Success           bool
	Type              string
	ProviderPaymentID string
	Data              any
}

func DispatchStripeEvent(event *stripe.Event) (PaymentStatus, error) {
	switch event.Type {

	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
			return PaymentStatus{Success: false},
				fmt.Errorf("failed to parse payment_intent.succeeded: %w", err)
		}

		return PaymentStatus{
			Success: true,
			Type:    "payment_intent.succeeded",
			Data:    paymentIntent,
		}, nil

	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		if err := json.Unmarshal(event.Data.Raw, &paymentMethod); err != nil {
			return PaymentStatus{Success: false},
				fmt.Errorf("failed to parse payment_method.attached: %w", err)
		}

		return PaymentStatus{
			Success: true,
			Type:    "payment_method.attached",
			Data:    paymentMethod,
		}, nil

	case "checkout.session.completed":
		var cs stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &cs); err != nil {
			return PaymentStatus{Success: false},
				fmt.Errorf("failed to parse checkout.session.completed: %w", err)
		}

		return PaymentStatus{
			Success:           true,
			Type:              "checkout.session.completed",
			ProviderPaymentID: cs.ID,
			Data:              cs,
		}, nil

	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
		return PaymentStatus{Success: false}, nil
	}
}
