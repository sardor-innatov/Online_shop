package cmd

import (
	"log"
	"online_shop/src/common/http/route"
	"online_shop/src/common/middleware"
	order_hanlder "online_shop/src/services/order_service/order/handler"
	migrate "online_shop/src/services/order_service/order_migrate"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Cmd(router route.Router, conn *pgxpool.Pool) {

	routerGroup := router.Group("/api/v1", middleware.RequireAuth)
	{
		order_hanlder.NewOrderHandler(conn, *routerGroup)
	}

	err := migrate.Migrate(conn)
	{
		if err != nil {
			log.Fatal("failed to migrate database : ", err.Error())
		}
	}

}
