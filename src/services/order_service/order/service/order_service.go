package service

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"online_shop/src/common/http/route"
	"online_shop/src/services/order_service/order/dto"
	"online_shop/src/services/order_service/order/model"
	"online_shop/src/services/order_service/order/repository"
	card_repository "online_shop/src/services/payment_service/card/repository"
	payment_dto "online_shop/src/services/payment_service/payment/dto"
	payment_service "online_shop/src/services/payment_service/payment/service"
	product_repository "online_shop/src/services/product_service/product/repository"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderService interface {
	Create(ctx route.Context) (*payment_dto.PaymentResponse, error)
}

type orderService struct {
	conn         *pgxpool.Pool
	card_repo    card_repository.CardRepository
	order_repo   repository.OrderRepository
	product_repo product_repository.ProductRepository
	payment      payment_service.PaymentService
}

func NewOrderService(conn *pgxpool.Pool) OrderService {
	return &orderService{
		conn:         conn,
		card_repo:    card_repository.NewCardRepository(conn),
		order_repo:   repository.NewOrderRepository(conn),
		product_repo: product_repository.NewProductRepository(conn),
		payment:      payment_service.NewPaymentService(),
	}
}

func (s *orderService) Create(ctx route.Context) (*payment_dto.PaymentResponse, error) {

	var order dto.OrderCreateDto
	err := ctx.Bind(&order)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid json"})
		}
	}

	val := ctx.Get("id")
	userId, ok := val.(int64)
	{
		if !ok {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid user id type in context"})
			return nil, errors.New("invalid user id type in context")
		}
	}
	// check if card registered

	card, err := s.card_repo.SelectOne(order.CardId)
	{
		if card == nil && err == nil {
			ctx.JSON(http.StatusNotFound, map[string]any{"error": "card not found"})
			return nil, sql.ErrNoRows
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]any{"error": "failed to SELECT card"})
			panic("failed to SELECT card")
		}
	}

	if card.UserId != userId {
		ctx.JSON(http.StatusNotFound, map[string]any{"error": "card not found"})
		return nil, nil // need error
	}

	// check products existense
	var ids []int64
	for _, item := range order.Items {
		ids = append(ids, item.ProductId)
	}

	ok, err = s.product_repo.CheckProducts(ids, ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return nil, err
		} else if !ok {
			ctx.JSON(http.StatusNotFound, map[string]any{"error": "items not found"})
			return nil, sql.ErrNoRows
		}
	}

	// decrease stack
	c, cancel := context.WithTimeout(context.Background(), time.Millisecond*15000)
	defer cancel()

	tx, err := s.conn.Begin(c)
	{
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return nil, err
		}
	}

	defer tx.Rollback(c)

	totalPrice, err := s.product_repo.UpdateStockItems(tx, order.Items, ctx)
	{
		if err != nil {
			tx.Rollback(c)
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "not enough items in stock"})
			return nil, err
		}
	}

	// create order

	orderModel := model.Order{
		UserId:     userId,
		TotalPrice: totalPrice,
		Status:     model.PENDING,
	}

	orderItemsWithPrice, err := s.product_repo.GetItemsWithPrices(order.Items, ctx)
	{
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
			return nil, err
		}
	}
	id, err := s.order_repo.Insert(tx, &orderModel, orderItemsWithPrice, ctx)
	{
		if err != nil {
			tx.Rollback(c)
			ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return nil, err
		}
	}

	// get payment

	paymentRequest := payment_dto.PaymentRequest{
		UserId:           strconv.FormatInt(userId, 10),
		Amount:           totalPrice,
		Currency:         "usd",
		StripeCustomerId: card.StripeCustomerId,
		PaymentMethodId:  card.StripeMethodId,
		OrderId:          strconv.FormatInt(id, 10),
	}

	resp, err := s.payment.PayOrder(paymentRequest)
	{
		if err != nil {
			tx.Rollback(c)
			ctx.JSON(http.StatusBadRequest, map[string]any{"error": "Payment failed: " + err.Error()})
			return nil, err
		}
	}

	//
	err = tx.Commit(c)
	{
		if err != nil {
			tx.Rollback(c)
			ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
			return nil, err
		}
	}
	return resp, nil
}
