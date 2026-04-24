package cmd

import (
	"context"
	"net/http"
	"online_shop/src/common/config"
	"online_shop/src/common/http/route"
	"online_shop/src/common/middleware"
	order "online_shop/src/services/order_service"
	product "online_shop/src/services/product_service"
	user "online_shop/src/services/user_service"
	payment "online_shop/src/services/payment_service"

	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Exec() {

	prjectEnv := configureEnv()

	println(prjectEnv.BaseURL)

	conn, err := pgxpool.New(context.Background(), "postgres://postgres:M067Ss19@localhost:5432/E_commerce?sslmode=disable")
	{
		if err != nil {
			println(err.Error())
		}
	}

	router := route.NewRoute(http.NewServeMux())

	router.Use(middleware.RateLimitMiddleware(config.GetRedis()))

	router.Handle("GET", "/swagger/", router.WrapHandler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	)))

	{
		user.Cmd(router, conn)
		product.Cmd(router, conn)
		order.Cmd(router, conn)
		payment.Cmd(router, conn)
	}

	createOrg(conn, router)

	router.Start()
}
