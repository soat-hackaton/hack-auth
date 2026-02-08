package main

import (
	"context"
	"log"

	"hack-auth/internal/api"    // <--- Importe o novo pacote api
	"hack-auth/internal/config"
	"hack-auth/internal/handler"
	"hack-auth/internal/repository/dynamo"
	"hack-auth/internal/service"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	// 1. Configuração
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// 2. Infraestrutura (AWS)
	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithRegion(cfg.AWSRegion))
	if err != nil {
		log.Fatalf("❌ Unable to load SDK config: %v", err)
	}
	dynamoClient := dynamodb.NewFromConfig(awsCfg)

	// 3. Injeção de Dependências (Wiring)
	userRepo := dynamo.NewUserRepository(dynamoClient, cfg.DynamoTableName)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	// 4. Setup do Servidor (Abstraído)
	// A main não sabe quais rotas existem, ela só pede um servidor pronto
	r := api.SetupRouter(authHandler)

	// 5. Start
	log.Printf("🚀 Auth Service running on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}