package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"online_shop/src/common/http/jwt"
	"online_shop/src/common/http/route"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RequireAuth(next route.HandlerFunc) route.HandlerFunc {
	return func(ctx route.Context) error {

		authHeader := ctx.Request.Header.Get("Authorization")

		tokenJSON, ok := jwt.ValidateToken(authHeader)
		{
			if !ok {
				return ctx.JSON(http.StatusUnauthorized, map[string]any{"error": "invalid token"})
			}
		}

		claims, err := jwt.GetClaims(tokenJSON)
		{
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
				return err
			}
		}

		// check if token is expired
		if float64(time.Now().Unix()) > float64(claims.ExpiresAt.Unix()) {
			log.Println("token expired")
			return ctx.JSON(http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		}

		ctx.Set("id", claims.Id)
		ctx.Set("role", claims.Role)
		return next(ctx)
	}
}

func CheckPermission(conn *pgxpool.Pool) route.Middleware {
	return func(next route.HandlerFunc) route.HandlerFunc {
		return func(ctx route.Context) error {
			log.Println("checking permissions")

			role := ctx.Get("role")
			path := ctx.Request.URL.Path
			method := ctx.Request.Method
			log.Println(role)

			permission := map[string][]string{
				path: {method},
			}

			jsonValue, _ := json.Marshal(permission)
			log.Println(jsonValue)
			ok, err := checkPermission(conn, role.(string), jsonValue)
			{
				if err != nil {
					log.Println(err.Error())
					return ctx.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
				}
				if !ok {
					log.Println(0)
					return ctx.JSON(http.StatusForbidden, map[string]any{"error": "you dont have permission to do this"})
				}
			}
			log.Println("next")
			return next(ctx)
		}
	}
}

func checkPermission(conn *pgxpool.Pool, role string, perrmission []byte) (bool, error) {

	query := `
	SELECT COUNT(id) FROM roles
	WHERE role_name = $1 AND permissions @> @$2 
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	var count int64
	err := conn.QueryRow(ctx, query, role, perrmission).Scan(&count)
	{
		if err != nil {
			return false, err
		} else if count == 0 {
			return false, nil
		}
	}

	return true, nil
}
