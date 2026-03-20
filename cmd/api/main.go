package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"flash-sale-order-system/internal/Infrastructure/idgen"
	"flash-sale-order-system/internal/Infrastructure/persistence/postgres"
	redisinfra "flash-sale-order-system/internal/Infrastructure/persistence/redis"
	infrarepo "flash-sale-order-system/internal/Infrastructure/persistence/repository"
	apporder "flash-sale-order-system/internal/application/order"
	appstock "flash-sale-order-system/internal/application/stock"
	httpserver "flash-sale-order-system/internal/interfaces/http"
	httporder "flash-sale-order-system/internal/interfaces/http/order"
	"flash-sale-order-system/internal/provider"
	"flash-sale-order-system/pkg/metrics"
	"flash-sale-order-system/pkg/rabbitmq"
)

func main() {
	// 1. Database
	db, err := postgres.NewDatabase(postgres.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "flashsale"),
		Password: getEnv("DB_PASSWORD", "flashsale123"),
		DBName:   getEnv("DB_NAME", "flashsale_db"),
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer postgres.CloseDatabase(db)

	// 2. ID Generator
	idGen, err := idgen.NewIDGenerator(1)
	if err != nil {
		log.Fatalf("failed to create id generator: %v", err)
	}

	// 3. Redis
	redisClient, err := redisinfra.NewClient(redisinfra.Config{
		Host: getEnv("REDIS_HOST", "localhost"),
		Port: getEnvInt("REDIS_PORT", 6379),
	})
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer redisinfra.CloseClient(redisClient)
	stockCache := redisinfra.NewStockCache(redisClient)

	// 4. RabbitMQ
	rmqURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	rmqConn, producerCh, err := rabbitmq.NewConnection(rmqURL)
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq: %v", err)
	}
	defer producerCh.Close()
	defer rmqConn.Close()

	consumerCh, err := rmqConn.Channel()
	if err != nil {
		log.Fatalf("failed to open consumer channel: %v", err)
	}
	defer consumerCh.Close()

	// 5. Order Producer
	producer := rabbitmq.NewOrderProducer(producerCh)

	// 6. Order Consumer — saves to PostgreSQL in the background
	saveOrderHandler := apporder.NewSaveOrderHandler(db, stockCache)
	consumer := rabbitmq.NewOrderConsumer(consumerCh)
	ctx := context.Background()
	if err := consumer.Start(ctx, func(msg rabbitmq.OrderMessage) error {
		return saveOrderHandler.Handle(ctx, msg)
	}); err != nil {
		log.Fatalf("failed to start order consumer: %v", err)
	}

	// 7. Stock Service — warm-up Redis from DB, then reconcile periodically
	productRepo := infrarepo.NewPostgresProductRepository(db)
	stockService := appstock.NewStockService(productRepo, stockCache)

	if err := stockService.WarmUp(ctx); err != nil {
		log.Printf("stock warm-up warning: %v", err)
	}
	go stockService.StartReconcile(ctx, 5*time.Minute)

	// 8. HTTP Handlers (via provider)
	productHandlers := provider.NewProductHandlers(db, idGen)
	placeOrderHandler := apporder.NewPlaceOrderHandler(stockCache, producer, idGen)

	handlers := &httpserver.Handlers{
		ProductCommand: productHandlers.Command,
		ProductQuery:   productHandlers.Query,
		Order:          httporder.NewHandler(placeOrderHandler),
	}

	// 9. Router
	router := httpserver.NewRouter(handlers)
	engine := router.Setup()

	// 10. Metrics Server
	metrics.Start(getEnv("METRICS_PORT", "9100"))

	// 11. Start Server
	port := getEnv("APP_PORT", "8080")
	log.Printf("Starting server on port %s...", port)
	if err := engine.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
