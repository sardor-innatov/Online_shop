package cmd

import (
	"log"
	"online_shop/src/common/http/route"
	category_handler "online_shop/src/services/product_service/category/handler"
	product_handler "online_shop/src/services/product_service/product/handler"
	migrate "online_shop/src/services/product_service/product_migrate"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Cmd(router route.Router, conn *pgxpool.Pool) {

	routerGroup := router.Group("/api/v1")
	{
		product_handler.NewProductHandler(conn, *routerGroup)
		category_handler.NewCategoryHandler(conn, *routerGroup)
	}

	err := migrate.Migrate(conn)
	{
		if err != nil {
			log.Fatal("failed to migrate database : ", err.Error())
		}
	}

}
