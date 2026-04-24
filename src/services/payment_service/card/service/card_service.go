package service

import (
	"fmt"
	"net/http"
	"online_shop/src/common/http/route"
	"online_shop/src/services/payment_service/card/dto"
	"online_shop/src/services/payment_service/card/model"
	"online_shop/src/services/payment_service/card/repository"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CardService interface {
	Create(ctx route.Context) (int64, error)
	CheckCard(cardId int64, ctx route.Context) bool
	UpdateBalance(ctx route.Context) error
}

type cardService struct {
	repo repository.CardRepository
}

func NewCardService(conn *pgxpool.Pool) CardService {
	return &cardService{
		repo: repository.NewCardRepository(conn),
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

	card := model.Card{
		UserId:        userid,
		Number:        dto.Number,
		Name:          dto.Name,
		PaymentSystem: dto.PaymentSystem,
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
