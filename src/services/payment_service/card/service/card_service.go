package service

import (
	"fmt"
	"net/http"
	"online_shop/src/common/config"
	"online_shop/src/common/http/route"
	"online_shop/src/services/payment_service/card/dto"
	"online_shop/src/services/payment_service/card/model"
	"online_shop/src/services/payment_service/card/repository"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stripe/stripe-go/v85"
	"github.com/stripe/stripe-go/v85/customer"
	"github.com/stripe/stripe-go/v85/paymentmethod"
	"github.com/stripe/stripe-go/v85/setupintent"
)

type CardService interface {
	Create(ctx route.Context) (int64, error)
	CheckCard(cardId int64, ctx route.Context) bool
	UpdateBalance(ctx route.Context) error
	SetupInitiate(ctx route.Context) (*dto.InitiateSetupResponseDto, error)
}

type cardService struct {
	repo            repository.CardRepository
	stripeSecretKey string
}

func NewCardService(conn *pgxpool.Pool) CardService {

	envProject := config.ProjectEnv()

	return &cardService{
		repo:            repository.NewCardRepository(conn),
		stripeSecretKey: envProject.StripeSecretKey,
	}
}

func (s *cardService) Create(ctx route.Context) (int64, error) {

	val := ctx.Get("id")
	userid, ok := val.(int64)
	{
		if !ok {
			return 0, ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid user id type in context"})
		}
	}

	var dto dto.CardCreateDto
	err := ctx.Bind(&dto)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid json"})
			return 0, err
		}
	}

	params := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(dto.StripeCustomerId),
	}

	_, err = paymentmethod.Attach(dto.StripeMethodId, params)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "failed to attach payment method"})
			return 0, err
		}
	}

	// 2. Делаем эту карту основной для клиента
	custParams := &stripe.CustomerParams{
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(dto.StripeMethodId),
		},
	}
	_, err = customer.Update(dto.StripeCustomerId, custParams)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "failed to update payment method"})
			return 0, err
		}
	}

	card := model.Card{
		UserId:           userid,
		StripeMethodId:   dto.StripeMethodId, // Например, "pm_1H6..."
		StripeCustomerId: dto.StripeCustomerId,
		Last4:            dto.Last4,
		PaymentSystem:    dto.PaymentSystem,
		Name:             dto.Name,
	}

	id, err := s.repo.Insert(card)
	{
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]any{"error": "failed to insert card"})
			fmt.Println(err)
			return 0, err
		}
	}

	return id, nil
}

func (s *cardService) CheckCard(cardId int64, ctx route.Context) bool {

	val := ctx.Get("id")
	userid, ok := val.(int64)
	{
		if !ok {
			ctx.JSON(http.StatusInternalServerError, map[string]any{"error": "failed to read user id"})
			panic("failed to read user id")
		}
	}

	card, err := s.repo.SelectOne(cardId)
	{
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]any{"error": "failed to SELECT card"})
			panic("failed to SELECT card")
		}
		if card == nil {
			ctx.JSON(http.StatusNotFound, map[string]any{"error": "card not found"})
			return false
		}

	}

	if card.UserId != userid {
		ctx.JSON(http.StatusNotFound, map[string]any{"error": "card not found"})
		return false
	}

	return true
}

func (s *cardService) UpdateBalance(ctx route.Context) error {

	idStr := ctx.PathParam("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid path param"})
			return err
		}
	}

	var dto dto.BalanceUpdateDto
	err = ctx.Bind(&dto)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid json"})
			return err
		}
	}

	err = s.repo.UpdateBalance(id, dto.Amaunt)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "lack of pounds"})
			return err
		}
	}

	return nil
}

func (s *cardService) SetupInitiate(ctx route.Context) (*dto.InitiateSetupResponseDto, error) {

	val := ctx.Get("id")
	userid, ok := val.(int64)
	{
		if !ok {
			return nil, ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid user id type in context"})
		}
	}

	var dtoReq dto.InitiateSetupRequestDto
	err := ctx.Bind(&dtoReq)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid json"})
			return nil, err
		}
	}

	stripe.Key = s.stripeSecretKey

	stripeCustomerId, err := s.repo.GetStripeCustomerId(userid, ctx)
	{
		if err != nil || stripeCustomerId == "" {

			params := &stripe.CustomerParams{
				Email: stripe.String(dtoReq.Email),
				Metadata: map[string]string{
					"user_id": fmt.Sprintf("%d", userid),
				},
			}
			newCust, err := customer.New(params)
			{
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, map[string]any{"error": "stripe customer creation failed"})
					return nil, err
				}
			}

			stripeCustomerId = newCust.ID
		}
	}

	params := &stripe.SetupIntentParams{
		Customer: stripe.String(stripeCustomerId),
		Usage:    stripe.String("off_session"),
	}

	si, err := setupintent.New(params)
	{
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return nil, err
		}
	}

	response := dto.InitiateSetupResponseDto{
		ClientSecret:     si.ClientSecret,
		StripeCustomerId: si.Customer.ID,
	}

	return &response, nil
}
