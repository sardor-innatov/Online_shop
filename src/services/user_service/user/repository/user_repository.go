package repository

import (
	"context"
	"fmt"
	"time"

	"online_shop/src/common/http/route"
	"online_shop/src/services/user_service/user/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Insert(user *model.User, c route.Context) (int64, error)
	Udpate(id int64, user *model.User, c route.Context) error
	Delete(id int64, c route.Context) error
	SelectOne(id int64, c route.Context) (*model.User, error)
	SelectByEmail(email string, c route.Context) (*model.User, error)
	SelectWithPage(skip int64, take int64, c route.Context) ([]model.User, error)
}

type userRepository struct {
	conn *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) UserRepository {
	return &userRepository{
		conn: conn,
	}
}

func (r *userRepository) Insert(user *model.User, c route.Context) (int64, error) {

	var query string

	query = "INSERT INTO \"users\" (first_name, last_name, email, password) "

	query += fmt.Sprintf("VALUES ('%s','%s','%s','%s')", user.FirstName, user.LastName, user.Email, user.Password)

	query += " RETURNING id"

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*200)
	defer cancel()

	var id int64
	err := r.conn.QueryRow(ctx, query).Scan(&id)
	{
		fmt.Println(query)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (r *userRepository) Udpate(id int64, user *model.User, c route.Context) error {

	var query string

	query = "UPDATE users SET "

	query = fmt.Sprintf(query+"first_name = '%s', last_name = '%s', email = '%s'", user.FirstName, user.LastName, user.Email)

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

func (r *userRepository) Delete(id int64, c route.Context) error {
	var query string

	query = "DELETE FROM users "

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

func (r *userRepository) SelectOne(id int64, c route.Context) (*model.User, error) {

	var query string

	query = "SELECT id, first_name, last_name, email FROM users "

	query = fmt.Sprintf(query+"WHERE id = %v", id)

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*200)
	defer cancel()

	var user model.User
	err := r.conn.QueryRow(ctx, query).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email)
	{
		fmt.Println(query)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (r *userRepository) SelectByEmail(email string, c route.Context) (*model.User, error) {

	var query string

	query = "SELECT id, first_name, last_name, email, password FROM users "

	query = fmt.Sprintf(query+"WHERE email = '%s'", email)

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*200)
	defer cancel()

	var user model.User
	err := r.conn.QueryRow(ctx, query).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	{
		fmt.Println(query)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (r *userRepository) SelectWithPage(skip int64, take int64, c route.Context) ([]model.User, error) {

	var query string

	query = "SELECT id, first_name, last_name, email FROM users "

	query = fmt.Sprintf(query+"LIMIT %v OFFSET %v", take, skip)

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*500)
	defer cancel()

	var users []model.User
	rows, err := r.conn.Query(ctx, query)
	{
		fmt.Println(query)
		if err != nil {
			return users, err
		}
	}

	defer rows.Close()

	for rows.Next() {

		var user model.User
		err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email)
		{
			if err != nil {
				return nil, err
			}
		}

		users = append(users, user)
	}

	return users, nil
}
