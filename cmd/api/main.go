package main

import (
	"log"
	"os"
	"strconv"

	"flash-sale-order-system/internal/Infrastructure/idgen"
	"flash-sale-order-system/internal/Infrastructure/persistence/postgres"
	persistence "flash-sale-order-system/internal/Infrastructure/persistence/repository"
	productapp "flash-sale-order-system/internal/application/product"
	httpserver "flash-sale-order-system/internal/interfaces/http"
	"flash-sale-order-system/internal/interfaces/http/handler"
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

	// 3. Repositories
	productRepo := persistence.NewPostgresProductRepository(db)
	pricingRepo := persistence.NewPostgresProductPricingRepository(db)

	// 4. Application Handlers
	createProductHandler := productapp.NewCreateProductHandler(db, idGen, productRepo, pricingRepo)

	// 5. HTTP Handlers
	productHandler := handler.NewProductHandler(createProductHandler)

	// 6. Router
	router := httpserver.NewRouter(productHandler)
	engine := router.Setup()

	// 7. Start Server
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
