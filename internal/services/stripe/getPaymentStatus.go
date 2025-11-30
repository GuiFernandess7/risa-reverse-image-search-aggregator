package stripe

import (
	"log"
	"os"

	"github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/checkout/session"
)

func GetPaymentStatus(pID string) (*stripe.CheckoutSession, error) {
	stripe.Key = os.Getenv("STRIPE_API_KEY")
	sess, err := session.Get(pID, nil)
	if err != nil {
		log.Printf("Error getting payment session: %v", err)
		return nil, err
	}
	return sess, nil
}
