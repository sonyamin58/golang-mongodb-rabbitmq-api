package middleware

import (
	"github.com/ibas/golib-api/internal/config"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTConfig(cfg *config.Config) echojwt.Config {
	return echojwt.Config{
		SigningKey: []byte(cfg.JWT.Secret),
		TokenLookup: "header:Authorization",
		AuthScheme: "Bearer",
		ContextKey:  "user",
		ErrorHandler: func(c echo.Context, err error) error {
			return echo.ErrUnauthorized
		},
	}
}
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc
