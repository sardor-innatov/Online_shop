package cmd

import (
	"context"
	"fmt"
	"online_shop/src/common/config"
	"online_shop/src/common/helper"
	"online_shop/src/common/http/route"
	role_model "online_shop/src/services/user_service/role/model"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func configureEnv() config.EnvProject {
	return config.ProjectEnv()
}

func createOrg(conn *pgxpool.Pool, r route.Router) {

	routers := make(helper.JsonObject)

	for _, route := range r.Routes() {
		// fmt.Println(route, i)
		if route.Method == "route_not_found" {

			continue
		}
		if existing, ok := routers[route.Path]; ok {
			routers[route.Path] = append(existing.([]string), route.Method)
		} else {
			routers[route.Path] = []string{route.Method}
		}
	}

	query := `
	DELETE FROM roles
	WHERE roles.id = 1 OR roles.id = 2`

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	_, err := conn.Exec(ctx, query)
	{
		if err != nil {
			panic(err.Error())
		}
	}

	fmt.Println(routers)

	roles := []role_model.Role{
		{
			Id:          1,
			RoleName:    "superadmin",
			Permissions: routers,
		},
		{
			Id:       2,
			RoleName: "user",
			Permissions: helper.JsonObject{
				// product
				"/api/v1/product": []string{"GET"},
				"/api/v1/product/{id}": []string{"GET"},
				// auth
				"/api/v1/auth":         []string{"POST"},
				"/api/v1/auth/login":   []string{"POST"},
				// category
				"/api/v1/category":      []string{"GET"},
				// order
				"/api/v1/order":      []string{"POST"},
				// card 
				"/api/v1/card":      []string{"POST"},
				"/api/v1/card/balance":      []string{"PUT"},
			},
		},
	}

	query = `
	INSERT INTO roles (id, role_name, permissions)
	VALUES ($1, $2, $3), ($4, $5, $6) `

	ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	_, err = conn.Exec(ctx, query, roles[0].Id, roles[0].RoleName, roles[0].Permissions, roles[1].Id, roles[1].RoleName, roles[1].Permissions)
	{
		if err != nil {
			panic(err.Error())
		}
	}

}
