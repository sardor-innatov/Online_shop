package repository

import (
	"context"
	"fmt"
	"online_shop/src/common/http/route"
	"online_shop/src/services/product_service/category/model"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository interface {
	Insert(category *model.Category, c route.Context) (int64, error)
	Update(id int64, category *model.Category, c route.Context) error
	Delete(id int64, c route.Context) error
	SelectOne(id int64, c route.Context) (*model.Category, error)
	SelectPage(take, skip int64, c route.Context) ([]model.Category, error)
}

type categoryRepository struct {
	conn *pgxpool.Pool
}

func NewCategoryRepository(conn *pgxpool.Pool) CategoryRepository {
	return &categoryRepository{
		conn: conn,
	}
}

func (r *categoryRepository) Insert(category *model.Category, c route.Context) (int64, error) {

	var query string

	query += "INSERT INTO categories (name) VALUES "

	query = fmt.Sprintf(query+"('%v')", category.Name)

	query += " RETURNING id"

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*200)
	defer cancel()

	var newId int64
	err := r.conn.QueryRow(ctx, query).Scan(&newId)
	{
		if err != nil {
			return 0, err
		}
	}

	return newId, nil
}

func (r *categoryRepository) Update(id int64, category *model.Category, c route.Context) error {

	var query string

	query += "UPDATE categories SET "

	query = fmt.Sprintf(query+"name = '%v' WHERE id = %v", category.Name, id)

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*200)
	defer cancel()

	_, err := r.conn.Exec(ctx, query)
	{
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *categoryRepository) Delete(id int64, c route.Context) error {

	var query string

	query += "DELETE FROM categories "

	query = fmt.Sprintf(query+"WHERE id = %v", id)

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*200)
	defer cancel()

	_, err := r.conn.Exec(ctx, query)
	{
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *categoryRepository) SelectOne(id int64, c route.Context) (*model.Category, error) {

	var query string

	query += "SELECT id, name FROM categories "

	query = fmt.Sprintf(query+"WHERE id = %v", id)

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*200)
	defer cancel()

	var category model.Category
	err := r.conn.QueryRow(ctx, query).Scan(&category.Id, &category.Name)
	{
		if err != nil {
			return nil, err
		}
	}

	return &category, nil
}

func (r *categoryRepository) SelectPage(take, skip int64, c route.Context) ([]model.Category, error) {

	var query string

	query += "SELECT id, name FROM categories "

	query = fmt.Sprintf(query+"LIMIT %v OFFSET %v", take, skip)

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*200)
	defer cancel()

	rows, err := r.conn.Query(ctx, query)
	{
		if err != nil {
			return nil, err
		}
	}

	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var category model.Category

		rows.Scan(&category.Id, &category.Name)

		categories = append(categories, category)
	}

	return categories, nil
}
