package main

import (
	"1dv027/aad/internal/config"
	"1dv027/aad/internal/router"
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// @title           DogAdoptionApp REST HATEOAS API
// @version         1.0
// @description     An API with centralized dog adoption information.

// @contact.name   Gustav Näslund, gn222gq
// @contact.email  gn222gq@student.lnu.se

// @host      https://cscloud7-113.lnu.se/dogadoption
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	envPath, err := filepath.Abs("../.env")
	if err != nil {
		log.Fatalf("Error getting absolute path to env")
	}
	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	context := context.Background()
	dbPool, err := pgxpool.New(context, os.Getenv("DATABASE_CONNECTION_STRING"))
	if err != nil {
		log.Print("Could not create db pool")
		os.Exit(1)
	}
	defer dbPool.Close()

	containerConfig := config.ContainerConfig{
		DatabaseConnector:     dbPool,
		CryptographySecretKey: os.Getenv("CRYPTO_KEY"),
		BasePath:              os.Getenv("BASE_PATH"),
		JwtSigningKey:         os.Getenv("JWT_SIGNING_KEY"),
	}

	container := config.SetupContainer(containerConfig)

	router := router.NewRouter(container)
	router.StartRouter()
}
