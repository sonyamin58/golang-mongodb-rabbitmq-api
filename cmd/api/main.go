package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ibas/golib-api/internal/config"
	"github.com/ibas/golib-api/internal/handler"
	"github.com/ibas/golib-api/internal/middleware"
	"github.com/ibas/golib-api/internal/repository"
	"github.com/ibas/golib-api/internal/service"
	"github.com/ibas/golib-api/pkg/database"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.InitOracle(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Redis
	redisClient, err := database.InitRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg)
	accountService := service.NewAccountService(accountRepo, userRepo)
	transactionService := service.NewTransactionService(transactionRepo, accountRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	accountHandler := handler.NewAccountHandler(accountService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())
	e.Use(echomw.CORS())

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Public routes
	public := e.Group("/api/v1")
	public.POST("/auth/register", authHandler.Register)
	public.POST("/auth/login", authHandler.Login)
	public.POST("/auth/refresh", authHandler.RefreshToken)

	// Rate limiting middleware
	if cfg.RateLimit.Enabled {
		public.Use(middleware.RateLimiter(redisClient, cfg))
	}

	// Protected routes
	protected := e.Group("/api/v1")
	protected.Use(echojwt.WithConfig(middleware.JWTConfig(cfg)))

	// Account routes
	protected.GET("/accounts", accountHandler.List)
	protected.GET("/accounts/:id", accountHandler.Get)
	protected.POST("/accounts", accountHandler.Create)
	protected.PUT("/accounts/:id", accountHandler.Update)
	protected.DELETE("/accounts/:id", accountHandler.Delete)

	// Transaction routes
	protected.GET("/transactions", transactionHandler.List)
	protected.GET("/transactions/:id", transactionHandler.Get)
	protected.POST("/transactions", transactionHandler.Create)
	protected.POST("/transactions/transfer", transactionHandler.Transfer)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	go func() {
		log.Printf("Starting server on %s", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	// Close connections
	if err := db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
	if err := redisClient.Close(); err != nil {
		log.Printf("Error closing Redis: %v", err)
	}

	log.Println("Server stopped")
}
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc
