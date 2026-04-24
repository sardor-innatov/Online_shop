package cmd

import (
	"log"
	"online_shop/src/common/http/route"
	card_handler "online_shop/src/services/payment_service/card/handler"
	migrate "online_shop/src/services/payment_service/payment_migrate"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Cmd(router route.Router, conn *pgxpool.Pool) {

	routerGroup := router.Group("/api/v1")
	{
		card_handler.NewCardHandler(*routerGroup, conn)
	}

	err := migrate.Migrate(conn)
	{
		if err != nil {
			log.Fatal("failed to migrate database : ", err.Error())
		}
	}

}
