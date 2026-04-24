package migrate

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(conn *pgxpool.Pool) error {

	var query string

	{
		query = cardTable(query)
	}

	ctx := context.Background()
	_, err := conn.Exec(ctx, query)
	{
		if err != nil {
			return err
		}
	}

	return nil
}
