package service

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"online_shop/src/common/config"
	"online_shop/src/common/http/route"
	"online_shop/src/services/payment_service/payment/dto"

	"github.com/stripe/stripe-go/v85"
	"github.com/stripe/stripe-go/v85/paymentintent"
	"github.com/stripe/stripe-go/v85/webhook"
)

type PaymentService interface {
	PayOrder(req dto.PaymentRequest) (*dto.PaymentResponse, error)
	CheckWebhook(ctx route.Context) error
}

type paymentService struct {
	stripeSecretKey string
}

func NewPaymentService() PaymentService {

	envProject := config.ProjectEnv()

	return &paymentService{
		stripeSecretKey: envProject.StripeSecretKey,
	}
}

func (s *paymentService) PayOrder(req dto.PaymentRequest) (*dto.PaymentResponse, error) {
	stripe.Key = s.stripeSecretKey

	params := &stripe.PaymentIntentParams{

		Amount:   stripe.Int64(int64(math.Round(req.Amount * 100))),
		Currency: stripe.String(req.Currency),
		Customer: stripe.String(req.StripeCustomerId),
		PaymentMethod: stripe.String(req.PaymentMethodId),

		Confirm:    stripe.Bool(true),
		OffSession: stripe.Bool(true),


		//ConfirmationMethod: stripe.String("automatic"),

		PaymentMethodOptions: &stripe.PaymentIntentPaymentMethodOptionsParams{
			Card: &stripe.PaymentIntentPaymentMethodOptionsCardParams{
				RequestThreeDSecure: stripe.String("automatic"),
			},
		},
	}

	params.AddMetadata("order_id", req.OrderId)

	pIntent, err := paymentintent.New(params)
	{
		if err != nil {
			return nil, err
		}
	}

	return &dto.PaymentResponse{
		UserId:        req.UserId,
		Status:        string(pIntent.Status),
		Success:       pIntent.Status == "succeeded",
		ClientSecret:  pIntent.ClientSecret,
		TransactionId: pIntent.ID,
	}, nil
}

func (s *paymentService) CheckWebhook(ctx route.Context) error {

	payload, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{"err": "failed to read payload"})
		return err
	}
	defer ctx.Request.Body.Close()
	endpointSecret := "whsec_0c45d137d90930db56877c06539ff24202a71b40d7d980ec8a5b7ecd1d081ce0"

	signature := ctx.Request.Header.Get("Stripe-Signature")

	event, err := webhook.ConstructEvent(payload, signature, endpointSecret)
	if err != nil {
		return err
	}

	switch event.Type {

	// Платеж успешно завершен
	case "payment_intent.succeeded":
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			return fmt.Errorf("error unmarshaling stripe data: %w", err)
		}

		// Достаем ваш ID заказа, который вы передавали в Metadata при создании платежа
		//orderID := pi.Metadata["order_id"]

		// ШАГ 3: Ваша бизнес-логика
		//return s.finalizeOrder(orderID, pi.ID)
	case "payment_intent.payment_failed":
		var pi stripe.PaymentIntent
		json.Unmarshal(event.Data.Raw, &pi)

		//orderID := pi.Metadata["order_id"]
		//return s.markOrderAsFailed(orderID)

	default:
		// Другие события можно просто игнорировать
		fmt.Printf("Unhandled event type: %s\n", event.Type)
	}

	println("webhook")
	return nil
}
