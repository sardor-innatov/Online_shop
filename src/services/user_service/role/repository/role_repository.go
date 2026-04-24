package repository

import (
	"context"
	"online_shop/src/common/http/route"
	"online_shop/src/services/user_service/role/model"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepository interface {
	Insert(role *model.Role, c route.Context) (int64, error)
	CheckPermission(role string, perrmission []byte, c route.Context) (bool, error)
	Delete(id int64, c route.Context) error
}

type roleRepository struct {
	conn *pgxpool.Pool
}

func NewRoleRepository(conn *pgxpool.Pool) RoleRepository {
	return &roleRepository{
		conn: conn,
	}
}

func (r *roleRepository) Insert(role *model.Role, c route.Context) (int64, error) {

	query := `
	INSERT INTO roles (role_name, permissions)
	VALUES ($1, $2)
	RETURNING id`

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*200)
	defer cancel()

	var id int64
	err := r.conn.QueryRow(ctx, query, role.RoleName, role.Permissions).Scan(&id)
	{
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (r *roleRepository) Delete(id int64, c route.Context) error {
	query := `
    DELETE FROM roles
	WHERE role_name = $1 
	`
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*300)
	defer cancel()

	_, err := r.conn.Exec(ctx, query, id)
	{
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *roleRepository) CheckPermission(role string, perrmission []byte, c route.Context) (bool, error) {

	query := `
	SELECT COUNT(id) FROM roles
	WHERE role_name = $1 AND permissions @> @$2 
	`

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*300)
	defer cancel()

	var count int64
	err := r.conn.QueryRow(ctx, query, role, perrmission).Scan(&count)
	{
		if err != nil {
			return false, err
		} else if count == 0 {
			return false, nil
		}
	}

	return true, nil
}
