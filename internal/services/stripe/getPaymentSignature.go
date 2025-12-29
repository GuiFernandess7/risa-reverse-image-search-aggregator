package stripe

import (
	"fmt"

	stripe "github.com/stripe/stripe-go/v80"
	"github.com/stripe/stripe-go/v80/webhook"
)

func GetPaymentEvent(body []byte, signature string, stripeSecret string) (*stripe.Event, error) {
	event, err := webhook.ConstructEventWithOptions(
		body,
		signature,
		stripeSecret,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid webhook signature: %w", err)
	}
	return &event, nil
}
