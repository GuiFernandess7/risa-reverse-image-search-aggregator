package stripe

import (
	"fmt"

	stripe "github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/webhook"
)

func GetPaymentEvent(body []byte, signature string, stripeSecret string) (*stripe.Event, error) {
	event, err := webhook.ConstructEvent(body, signature, stripeSecret)
	if err != nil {
		return nil, fmt.Errorf("invalid webhook signature: %w", err)
	}
	return &event, nil
}
