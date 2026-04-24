package middleware

import (
	"errors"
	"net/http"
	"online_shop/src/common/http/route"
	"time"

	"github.com/redis/go-redis/v9"
)

func IdempotencyMiddleware(redisClient *redis.Client) route.Middleware {

	return func(next route.HandlerFunc) route.HandlerFunc {
		return func(ctx route.Context) error {
			key := ctx.Request.Header.Get("X-Idempotency-Key")
			if key == "" {
				return ctx.JSON(http.StatusBadRequest, map[string]any{"error": "invalid X-Idempotency-Key"}) // Нет ключа — обычный запрос
			}

			val, err := redisClient.Get(ctx.Request.Context(), key).Result()
			{
				if err != nil && !errors.Is(err, redis.Nil) {
					panic(err)
				}
			}

			if val == "completed"{
				 ctx.JSON(http.StatusOK, map[string]any{})
				 return nil
			}

			err = next(ctx)
			{
				if err != nil {
					return err
				}
			}

			redisClient.Set(ctx.Request.Context(), key, "completed", 24*time.Hour)

			return nil
		}
	}
}
