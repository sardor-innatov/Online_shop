package repository

import (
	"context"
	"fmt"
	"online_shop/src/common/http/route"
	order_dto "online_shop/src/services/order_service/order/dto"
	order_model "online_shop/src/services/order_service/order/model"
	"online_shop/src/services/product_service/product/dto"
	"online_shop/src/services/product_service/product/model"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository interface {
	Insert(product *model.Product, c route.Context) (int64, error)
	Update(id int64, product *model.Product, c route.Context) error
	UpdateStock(id int64, quantity int, c route.Context) error
	Delete(id int64, c route.Context) error
	SelectOne(id int64, c route.Context) (*model.Product, error)
	SelectOneWithCategoryName(id int64, c route.Context) (*dto.ProductGetDto, error)
	SelectWithPage(skip, take, categoryId int64, priceMin, priceMax float32, c route.Context) ([]dto.ProductGetDto, error)
	CheckProducts(ids []int64, c route.Context) (bool, error)
	UpdateStockItems(tx pgx.Tx, items []order_dto.OrderItemCreateDto, c route.Context) (float64, error)
	GetItemsWithPrices(items []order_dto.OrderItemCreateDto, c route.Context) ([]order_model.OrderItem, error)
}

type productRepository struct {
	conn *pgxpool.Pool
}

func NewProductRepository(conn *pgxpool.Pool) ProductRepository {
	return &productRepository{
		conn: conn,
	}
}

func (r *productRepository) Insert(product *model.Product, c route.Context) (int64, error) {

	var query string

	query = "INSERT INTO \"products\" (name, description, price, stock, category_id) "

	query += fmt.Sprintf("VALUES ('%s','%s', %v, %v, %v)",
		product.Name, product.Description, product.Price,
		product.Stock, product.CategoryId)

	query += " RETURNING id"

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*200)
	defer cancel()

	var newId int64
	err := r.conn.QueryRow(ctx, query).Scan(&newId)
	{
		fmt.Println(query)
		if err != nil {
			return 0, err
		}
	}

	return newId, nil
}

func (r *productRepository) Update(id int64, product *model.Product, c route.Context) error {

	var query string

	query = "UPDATE products SET "

	query = fmt.Sprintf(query+"name = '%s', description = '%s', price = '%v', stock = %v, category_id = %v, updated_at = Now()",
		product.Name, product.Description, product.Price,
		product.Stock, product.CategoryId)

	query = fmt.Sprintf(query+"WHERE id = %v", id)

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*200)
	defer cancel()
	_, err := r.conn.Exec(ctx, query)
	{
		fmt.Println(query)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *productRepository) UpdateStock(id int64, quantity int, c route.Context) error {

	var query string

	query = "UPDATE products SET "

	query = fmt.Sprintf(query+"stock = stock + %v, updated_at = Now() ", quantity)

	query = fmt.Sprintf(query+"WHERE id = %v", id)

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*200)
	defer cancel()
	_, err := r.conn.Exec(ctx, query)
	{
		fmt.Println(query)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *productRepository) Delete(id int64, c route.Context) error {
	var query string

	query = "DELETE FROM products "

	query = fmt.Sprintf(query+"WHERE id = %v", id)

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*200)
	defer cancel()
	_, err := r.conn.Exec(ctx, query)
	{
		fmt.Println(query)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *productRepository) SelectOne(id int64, c route.Context) (*model.Product, error) {

	var query string

	query = "SELECT name, description, price, stock, category_id, created_at, updated_at FROM products "

	query = fmt.Sprintf(query+"WHERE id = %v", id)

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*200)
	defer cancel()

	var product model.Product
	err := r.conn.QueryRow(ctx, query).Scan(&product.Name, &product.Description, &product.Price, &product.Stock, &product.CategoryId, &product.CreatedAt, &product.UpdatedAt)
	{
		fmt.Println(query)
		if err != nil {
			return nil, err
		}
	}

	return &product, nil
}

func (r *productRepository) SelectOneWithCategoryName(id int64, c route.Context) (*dto.ProductGetDto, error) {

	var query string

	query = `SELECT products.name, description, price, stock, categories.name, created_at, updated_at FROM products 
			JOIN categories on products.category_id = categories.id `

	query = fmt.Sprintf(query+"WHERE products.id = %v ", id)

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*500)
	defer cancel()

	var product dto.ProductGetDto
	err := r.conn.QueryRow(ctx, query).Scan(&product.Name, &product.Description, &product.Price, &product.Stock, &product.CategoryName, &product.CreatedAt, &product.UpdatedAt)
	{
		fmt.Println(query)
		if err != nil {
			return nil, err
		}
	}

	return &product, nil
}

func (r *productRepository) SelectWithPage(skip, take, categoryId int64, priceMin, priceMax float32, c route.Context) ([]dto.ProductGetDto, error) {

	var query string

	query = `SELECT products.name, description, price, stock, categories.name, created_at, updated_at FROM products 
			JOIN categories on products.category_id = categories.id `

	if categoryId != 0 || priceMin != 0 || priceMax != 0 {
		query += " WHERE "
		if categoryId != 0 {
			query = fmt.Sprintf(query+"products.category_id = %v", categoryId)
		}
		if priceMin != 0 {
			query = fmt.Sprintf(query+"products.price >= %v", priceMin)
		}
		if priceMax != 0 {
			query = fmt.Sprintf(query+"products.price <= %v", priceMax)
		}
	}

	query = fmt.Sprintf(query+"LIMIT %v OFFSET %v", take, skip)

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*500)
	defer cancel()

	var products []dto.ProductGetDto
	rows, err := r.conn.Query(ctx, query)
	{
		fmt.Println(query)
		if err != nil {
			return products, err
		}
	}

	defer rows.Close()

	for rows.Next() {

		var product dto.ProductGetDto
		err := rows.Scan(&product.Name, &product.Description, &product.Price, &product.Stock, &product.CategoryName, &product.CreatedAt, &product.UpdatedAt)
		{
			if err != nil {
				return nil, err
			}
		}

		products = append(products, product)
	}

	return products, nil
}

func (r *productRepository) CheckProducts(ids []int64, c route.Context) (bool, error) {

	var allExists bool

	query := `SELECT ARRAY(SELECT id FROM users WHERE id = ANY($1)) @> $1::bigint[]`

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*200)
	defer cancel()

	err := r.conn.QueryRow(ctx, query, ids).Scan(&allExists)
	{
		if err != nil {
			return false, nil
		}
	}

	return allExists, nil
}

// func (r *productRepository) UpdateStockItems(tx pgx.Tx, items []order_dto.OrderItemCreateDto) error {

// 	var query string
// 	for _, item := range items {

// 		query += "UPDATE products SET "

// 		query = fmt.Sprintf(query+"stock = stock + %v, updated_at = '%v'", item.Quantity, time.Now())

// 		query = fmt.Sprintf(query+"WHERE id = %v", item.ProductId)

// 		query += "; "
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
// 	defer cancel()
// 	_, err := tx.Exec(ctx, query)
// 	{
// 		fmt.Println(query)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func (r *productRepository) UpdateStockItems(tx pgx.Tx, items []order_dto.OrderItemCreateDto, c route.Context) (float64, error) {
	query := `
	WITH order_data AS (
    SELECT 
        unnest($1::int[]) as item_id, 
        unnest($2::int[]) as requested_qty
	),
	calculated AS (
	
	    SELECT 
	        i.id,
	        i.price,
	        d.requested_qty,
	        (i.price * d.requested_qty) as subtotal
	    FROM products i
	    JOIN order_data d ON i.id = d.item_id
	    --WHERE i.stock >= d.requested_qty
	    FOR UPDATE 
	),
	updated AS (
	    UPDATE products i SET 
		stock = i.stock - c.requested_qty,
		updated_at = CURRENT_TIMESTAMP
	    FROM calculated c
	    WHERE i.id = c.id
	      AND (SELECT count(*) FROM calculated) = array_length($1::int[], 1)
	    RETURNING c.subtotal
	)
	
	SELECT COALESCE(sum(subtotal), 0) as total_price FROM updated;
`
	var ids []int64
	var qty []int
	for _, item := range items {
		ids = append(ids, item.ProductId)
		qty = append(qty, item.Quantity)
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*1000)
	defer cancel()

	var totalPrice float64
	err := tx.QueryRow(ctx, query, ids, qty).Scan(&totalPrice)
	{
		if err != nil {
			return 0, err
		}
	}

	return totalPrice, nil
}

func (r *productRepository) GetItemsWithPrices(items []order_dto.OrderItemCreateDto, c route.Context) ([]order_model.OrderItem, error) {

	query := "SELECT id, price FROM products WHERE id = ANY($1)"

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*200)
	defer cancel()

	var ids []int64
	for _, item := range items {
		ids = append(ids, item.ProductId)
	}

	rows, err := r.conn.Query(ctx, query, ids)
	{
		if err != nil {
			return nil, err
		}
	}

	defer rows.Close()

	var orderItems []order_model.OrderItem
	for _, item := range items {
		if rows.Next() {

			var orderItem order_model.OrderItem

			orderItem.Quantity = item.Quantity

			err := rows.Scan(&orderItem.ProductId, &orderItem.Price)
			{
				if err != nil {
					return nil, err
				}
			}

			orderItems = append(orderItems, orderItem)
		}
	}

	return orderItems, nil
}
