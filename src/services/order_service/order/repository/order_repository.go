package repository

import (
	"context"
	"fmt"
	"online_shop/src/common/http/route"
	"online_shop/src/services/order_service/order/model"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	Insert(tx pgx.Tx, order *model.Order, items []model.OrderItem , c route.Context) (int64, error)
}

type orderRepository struct {
	conn *pgxpool.Pool
}

func NewOrderRepository(conn *pgxpool.Pool) OrderRepository {
	return &orderRepository{
		conn: conn,
	}
}

func (r *orderRepository) Insert(tx pgx.Tx, order *model.Order, items []model.OrderItem, c route.Context) (int64, error) {

	var query string

	query += "INSERT INTO orders (user_id, total_price, status) "

	query = fmt.Sprintf(query+"VALUES (%v, %v, '%s')", order.UserId, order.TotalPrice, order.Status)

	query += " RETURNING id;"

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*200)
	defer cancel()

	var newId int64
	err := tx.QueryRow(ctx, query).Scan(&newId)
	{
		if err != nil {
			return 0, err
		}
	}

	query = "INSERT INTO order_items (order_id, product_id, quantity, price) VALUES "

	for i, item := range items {
		query = fmt.Sprintf(query+" (%v, %v, %v, %v)", newId, item.ProductId, item.Quantity, item.Price)
		if i+1 != len(items) {
			query += ","
		}
	}

	_, err = tx.Exec(ctx, query)
	{
		if err != nil {
			return 0, err
		}
	}

	return newId, nil
}
