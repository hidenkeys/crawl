package services

import (
	"crawl/repositories"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/webhook"
	"log"
	"os"
)

type PaymentService interface {
	CreateCheckoutSession(req CreateCheckoutSessionRequest) (*CreateCheckoutSessionResponse, error)
	HandleWebhook(payload []byte, signature string) error
	fulfillPurchase(checkoutSession stripe.CheckoutSession) error
}

type paymentService struct {
	albumPurchaseRepo repositories.IAlbumPurchaseRepository
	songPurchaseRepo  repositories.ISongPurchaseRepository
	albumRepo         repositories.IAlbumRepository
	songRepo          repositories.ISongRepository
}

func NewPaymentService(
	albumPurchaseRepo repositories.IAlbumPurchaseRepository,
	songPurchaseRepo repositories.ISongPurchaseRepository,
	albumRepo repositories.IAlbumRepository,
	songRepo repositories.ISongRepository,
) *paymentService {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	return &paymentService{
		albumPurchaseRepo: albumPurchaseRepo,
		songPurchaseRepo:  songPurchaseRepo,
		albumRepo:         albumRepo,
		songRepo:          songRepo,
	}
}

type CreateCheckoutSessionRequest struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	Price    int64  `json:"price"`
	ItemType string `json:"item_type"` // "song" or "album"
	ItemID   string `json:"item_id"`
	UserID   string `json:"user_id"`
}

type CreateCheckoutSessionResponse struct {
	StripeSessionID string
	MetaData        map[string]string
	PaymentUrl      string
}

func (s *paymentService) CreateCheckoutSession(req CreateCheckoutSessionRequest) (*CreateCheckoutSessionResponse, error) {
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(stripe.CurrencyUSD)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(req.Name),
					},
					UnitAmount: stripe.Int64(req.Price),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(os.Getenv("SUCCESS_URL")),
		CancelURL:  stripe.String(os.Getenv("FAILURE_URL")),
		Metadata: map[string]string{
			"user_id":   req.UserID,
			"trans_id":  req.ID,
			"item_id":   req.ItemID,
			"item_type": req.ItemType,
		},
	}

	sess, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe checkout session: %w", err)
	}

	return &CreateCheckoutSessionResponse{
		StripeSessionID: sess.ID,
		MetaData: map[string]string{
			"user_id":   req.UserID,
			"trans_id":  req.ID,
			"item_id":   req.ItemID,
			"item_type": req.ItemType,
		},
		PaymentUrl: sess.URL,
	}, nil
}

func (s *paymentService) HandleWebhook(payload []byte, signature string) error {
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	event, err := webhook.ConstructEvent(payload, signature, endpointSecret)
	if err != nil {
		log.Printf("Failed to construct webhook event: %s", err)
		return fmt.Errorf("failed to construct webhook event: %w", err)
	}

	// Handle the event
	switch event.Type {
	case "checkout.session.completed", "checkout.session.async_payment_succeeded":
		var checkoutSession stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &checkoutSession)
		if err != nil {
			log.Printf("failed to unmarshal checkout session: %w", err)
			return fmt.Errorf("failed to unmarshal checkout session: %w", err)
		}

		err = s.fulfillPurchase(checkoutSession)
		if err != nil {
			log.Printf("Failed to fulfil purchase")
			return err
		}

		log.Printf("Checkout session completed for user %s, item %s (%s)",
			checkoutSession.Metadata["user_id"],
			checkoutSession.Metadata["item_id"],
			checkoutSession.Metadata["item_type"],
		)
	case "checkout.session.async_payment_failed":
		var checkoutSession stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &checkoutSession)
		if err != nil {
			return fmt.Errorf("failed to unmarshal checkout session: %w", err)
		}
		log.Printf("Checkout session payment failed for user %s, item %s (%s)",
			checkoutSession.Metadata["user_id"],
			checkoutSession.Metadata["item_id"],
			checkoutSession.Metadata["item_type"],
		)
	case "checkout.session.expired":
		var checkoutSession stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &checkoutSession)
		if err != nil {
			return fmt.Errorf("failed to unmarshal checkout session: %w", err)
		}
		log.Printf("Checkout session expired for user %s, item %s (%s)",
			checkoutSession.Metadata["user_id"],
			checkoutSession.Metadata["item_id"],
			checkoutSession.Metadata["item_type"],
		)
	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}

	return nil
}

func (s *paymentService) fulfillPurchase(checkoutSession stripe.CheckoutSession) error {
	itemType := checkoutSession.Metadata["item_type"]
	_, err := uuid.Parse(checkoutSession.Metadata["user_id"])
	if err != nil {
		log.Printf("Failed to parse user ID")
		return fmt.Errorf("failed to parse user ID: %w", err)
	}
	itemID, err := uuid.Parse(checkoutSession.Metadata["item_id"])
	if err != nil {
		log.Printf("Failed to parse item id")
		return fmt.Errorf("failed to parse item ID: %w", err)
	}

	purchaseID, err := uuid.Parse(checkoutSession.Metadata["trans_id"])
	if err != nil {
		log.Printf("Failed to parse purchase id")
		return fmt.Errorf("failed to parse purchase ID: %w", err)
	}

	switch itemType {
	case "song":
		purchase, err := s.songPurchaseRepo.GetByID(purchaseID)
		if err != nil {
			log.Printf("Failed to find purchase record")
			return fmt.Errorf("failed to find purchase record: %w", err)
		}

		_, err = s.songRepo.GetByID(itemID)
		if err != nil {
			return fmt.Errorf("failed to get song: %w", err)
		}
		purchase.PaymentStatus = "completed"
		//purchase.StripeTransactionID = checkoutSession.PaymentIntent.ID
		_, err = s.songPurchaseRepo.Update(purchase)
		if err != nil {
			return fmt.Errorf("failed to update song purchase: %w", err)
		}
	case "album":
		purchase, err := s.albumPurchaseRepo.GetByID(purchaseID)
		if err != nil {
			log.Printf("Failed to find purchase record")
			return fmt.Errorf("failed to find purchase record: %w", err)
		}

		_, err = s.albumRepo.GetByID(itemID)
		if err != nil {
			return fmt.Errorf("failed to get album: %w", err)
		}
		purchase.PaymentStatus = "completed"
		//purchase.StripeTransactionID = checkoutSession.PaymentIntent.ID
		_, err = s.albumPurchaseRepo.Update(purchase)
		if err != nil {
			return fmt.Errorf("failed to update album purchase: %w", err)
		}
	default:
		return fmt.Errorf("unknown item type: %s", itemType)
	}

	return nil
}
