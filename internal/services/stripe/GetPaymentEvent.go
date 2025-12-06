package stripe

import (
	"encoding/json"
	"fmt"
	"os"

	stripe "github.com/stripe/stripe-go/v80"
)

type PaymentStatus struct {
	Success bool
	Data    any
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
			Data:    paymentMethod,
		}, nil

	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
		return PaymentStatus{Success: false}, nil
	}
}
