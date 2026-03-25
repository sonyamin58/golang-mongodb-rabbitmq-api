package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ibas/golib-api/internal/config"
	"github.com/ibas/golib-api/internal/machinery"
	"github.com/ibas/golib-api/pkg/database"
)

// WorkerConfig holds worker configuration
type WorkerConfig struct {
	Concurrency int
	BrokerURL   string
	ResultURL   string
}

func main() {
	log.Println("Starting Machinery Worker...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Override with defaults if not set in config
	workerCfg := WorkerConfig{
		Concurrency: cfg.Machinery.Concurrency,
		BrokerURL:   cfg.Machinery.Broker,
		ResultURL:   cfg.Machinery.ResultBackend,
	}

	// Set defaults
	if workerCfg.Concurrency == 0 {
		workerCfg.Concurrency = 4
	}
	if workerCfg.BrokerURL == "" {
		workerCfg.BrokerURL = "redis://localhost:6379/0"
	}
	if workerCfg.ResultURL == "" {
		workerCfg.ResultURL = "redis://localhost:6379/1"
	}

	log.Printf("Worker Config: Concurrency=%d, Broker=%s, ResultBackend=%s",
		workerCfg.Concurrency, workerCfg.BrokerURL, workerCfg.ResultURL)

	// Initialize database connection
	db, err := database.InitOracle(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected")

	// Set task context for worker
	machinery.SetTaskContext(&machinery.TaskContext{DB: db})

	// Create Machinery server
	server, err := machinery.NewServer(workerCfg.BrokerURL, workerCfg.ResultURL)
	if err != nil {
		log.Fatalf("Failed to create Machinery server: %v", err)
	}

	// Register tasks
	if err := machinery.RegisterTasks(server); err != nil {
		log.Fatalf("Failed to register tasks: %v", err)
	}
	log.Println("Tasks registered")

	// Create worker
	worker := server.NewWorker("machinery_worker", workerCfg.Concurrency)

	// Error handler
	worker.SetErrorHandler(func(err error) {
		log.Printf("Worker error: %v", err)
	})

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down worker...")
		worker.Quit()
	}()

	// Start worker
	log.Printf("Starting worker with concurrency: %d", workerCfg.Concurrency)
	if err := worker.Launch(); err != nil {
		log.Fatalf("Worker failed: %v", err)
	}

	// Cleanup
	log.Println("Closing database connections...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	log.Println("Worker stopped")
}
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc
