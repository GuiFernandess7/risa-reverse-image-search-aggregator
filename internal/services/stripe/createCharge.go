package stripe

import (
	"fmt"
	"os"

	"github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/checkout/session"
)

func CreateCheckoutSession(userID int64, creditAmount int64, priceCents int64) (*stripe.CheckoutSession, error) {
	stripe.Key = os.Getenv("STRIPE_API_KEY")

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String("payment"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Quantity: stripe.Int64(1),
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(fmt.Sprintf("%d Credits", creditAmount)),
					},
					UnitAmount: stripe.Int64(priceCents),
				},
			},
		},
		SuccessURL: stripe.String("https://mockfront.com/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String("https://mockfront.com/cancel?session_id={CHECKOUT_SESSION_ID}"),
		Metadata: map[string]string{
			"user_id":       fmt.Sprint(userID),
			"credit_amount": fmt.Sprint(creditAmount),
		},
	}

	return session.New(params)
}
