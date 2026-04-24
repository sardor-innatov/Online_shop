package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"online_shop/src/services/payment_service/card/model"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CardRepository interface {
	Insert(card model.Card) (int64, error)
	SelectOne(cardId int64) (*model.Card, error)
	PayOrder(card model.Card, totalPrice float32) bool
	UpdateBalance(id int64, amaunt float64) error
}

type cardRepository struct {
	conn *pgxpool.Pool
}

func NewCardRepository(conn *pgxpool.Pool) CardRepository {
	return &cardRepository{
		conn: conn,
	}
}

func (r *cardRepository) Insert(card model.Card) (int64, error) {

	query := `
	INSERT INTO cards (user_id, name, number, payment_system)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	var id int64
	err := r.conn.QueryRow(ctx, query, card.UserId, card.Name, card.Number, card.PaymentSystem).Scan(&id)
	{
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (r *cardRepository) SelectOne(cardId int64) (*model.Card, error) {

	query := `
	SELECT id, user_id, number, name, payment_system, balance FROM cards
	WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	var card model.Card
	err := r.conn.QueryRow(ctx, query, cardId).Scan(&card.Id, &card.UserId, &card.Number, &card.Name, &card.PaymentSystem, &card.Balance)
	{
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil // if card not found not returnning error
			}
			fmt.Println(sql.ErrNoRows)
			fmt.Println(err)
			return nil, err
		}
	}

	return &card, nil
}

func (r *cardRepository) PayOrder(card model.Card, totalPrice float32) bool {

	// query := `
	// SELECT total_price FROM orders
	// WHERE id = $1`

	// ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	// defer cancel()

	// var totalPrice float64
	// err := r.conn.QueryRow(ctx, query, orderId).Scan(&totalPrice)
	// {
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// }

	if (card.Balance - totalPrice) < 0 {
		return false
	}

	query := `
	UPDATE cards SET
	balance = $2
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	_, err := r.conn.Exec(ctx, query, card.Id, (card.Balance - totalPrice))
	{
		if err != nil {
			panic(err.Error())
		}
	}

	return true
}

func (r *cardRepository) UpdateBalance(id int64, amaunt float64) error {

	query := `
	UPDATE cards SET
	balance = balance + $2
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	_, err := r.conn.Exec(ctx, query, id, amaunt)
	{
		if err != nil {
			return err
		}
	}

	return nil
}
