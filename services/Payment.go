package services

/*
import (
	"context"
	"crawl/models"
	"fmt"
	"github.com/stripe/stripe-go/v82"
	"os"
)

type Service struct {
	stripeKey     string
	webhookSecret string
}

func NewStripeService() *Service {
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	stripe.Key = stripeKey

	return &Service{
		stripeKey:     stripeKey,
		webhookSecret: webhookSecret,
	}
}

type PaymentItem struct {
	ID       string
	Name     string
	Price    int64  // in cents
	Currency string // e.g., "usd"
}

func (s *Service) CreatePaymentLink(ctx context.Context, item PaymentItem, idempotencyKey string) (*models.Transaction, string, error) {
	params := &stripe.PaymentLinkParams{
		LineItems: []*stripe.PaymentLinkLineItemParams{
			{
				PriceData: &stripe.PaymentLinkLineItemPriceDataParams{
					Currency:   stripe.String(item.Currency),
					UnitAmount: stripe.Int64(item.Price),
					ProductData: &stripe.PaymentLinkLineItemPriceDataProductDataParams{
						Name: stripe.String(item.Name),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		Params: stripe.Params{
			Context:        ctx,
			IdempotencyKey: stripe.String(idempotencyKey),
		},
		AfterCompletion: &stripe.PaymentLinkAfterCompletionParams{
			Type: stripe.String("redirect"),
			Redirect: &stripe.PaymentLinkAfterCompletionRedirectParams{
				URL: stripe.String("your-app.com/payment-complete"), // Your frontend URL
			},
		},
	}

	pl, err := paymentlink.New(params)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create payment link: %w", err)
	}

	transaction := &models.Transaction{
		StripeID:       pl.ID,
		Amount:         item.Price,
		Currency:       item.Currency,
		Status:         "pending",
		IdempotencyKey: idempotencyKey,
		Metadata: models.JSON{
			"item_id":   item.ID,
			"item_name": item.Name,
		},
	}

	return transaction, pl.URL, nil
}

func (s *Service) ConstructEvent(payload []byte, signature string) (stripe.Event, error) {
	return webhook.ConstructEvent(payload, signature, s.webhookSecret)
}
*/
