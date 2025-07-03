package services

import (
	"context"
	"crawl/models"
	"github.com/google/uuid"

	"net/http"
	"os"

	_ "github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

func InitStripe() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
}

func GeneratePaymentLink(c context.Context, song models.Song, artistName string, userID uuid.UUID) error {
	//songID := c.Param("songId")
	//price := c.QueryParam("price") // Price in cents (e.g., 199 = $1.99)

	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("ngn"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Song Purchase for " + song.Title + " by " + artistName + " " + song.ID.String()),
					},
					UnitAmount: stripe.Int64(int64(song.Price)),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String("https://yourdomain.com/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String("https://yourdomain.com/cancel"),
		Metadata: map[string]string{
			"song_id": song.ID.String(),
			"user_id": userID.String(),
		},
	}

	s, err := session.New(params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"url": s.URL})
}
