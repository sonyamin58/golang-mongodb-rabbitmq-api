package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ibas/golib-api/internal/config"
	"github.com/ibas/golib-api/internal/response"
	"github.com/go-redis/redis/v9"
	"github.com/labstack/echo/v4"
)

func RateLimiter(redisClient *redis.Client, cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get client IP
			ip := c.RealIP()
			key := fmt.Sprintf("rate_limit:%s", ip)

			// Get current request count
			ctx := context.Background()
			count, err := redisClient.Get(ctx, key).Int()
			if err != nil && err != redis.Nil {
				// If Redis error, allow request through
				return next(c)
			}

			// Check if rate limit exceeded
			if count >= cfg.RateLimit.RequestsPerMinute {
				return c.JSON(http.StatusTooManyRequests, response.Error("Rate limit exceeded"))
			}

			// Increment counter
			pipe := redisClient.Pipeline()
			pipe.Incr(ctx, key)
			if count == 0 {
				pipe.Expire(ctx, key, time.Minute)
			}
			_, err = pipe.Exec(ctx)
			if err != nil {
				// If Redis error, allow request through
				return next(c)
			}

			return next(c)
		}
	}
}
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc
