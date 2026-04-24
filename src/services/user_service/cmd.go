package cmd

import (
	"log"
	"online_shop/src/common/http/route"
	auth_handler "online_shop/src/services/user_service/auth/handler"
	role_handler "online_shop/src/services/user_service/role/handler"
	user_handler "online_shop/src/services/user_service/user/handler"
	migrate "online_shop/src/services/user_service/user_migrate"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Cmd(router route.Router, conn *pgxpool.Pool) {

	routerGroup := router.Group("/api/v1")
	{
		user_handler.NewUserHandler(conn, *routerGroup)
		auth_handler.NewAuthHandler(conn, *routerGroup)
		role_handler.NewRoleHandler(conn, *routerGroup)
	}

	err := migrate.Migrate(conn)
	{
		if err != nil {
			log.Fatal("failed to migrate database : ", err.Error())
		}
	}

}
