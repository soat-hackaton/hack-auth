package main

import (
	"context"
	"log/slog"

	"hack-auth/internal/api"
	"hack-auth/internal/config"
	"hack-auth/internal/handler"
	"hack-auth/internal/repository/dynamo"
	"hack-auth/internal/service"
	"hack-auth/internal/utils/logger"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	// Setup do Logger
	logger.Init()
	slog.Info("Initializing Hack Auth Service...")

	// Configuração e Setup seguindo a ordem de dependências:
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		panic(err)
	}

	// Infraestrutura (AWS)
	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithRegion(cfg.AWSRegion))
	if err != nil {
		slog.Error("Unable to load SDK config", "error", err)
		panic(err)
	}
	dynamoClient := dynamodb.NewFromConfig(awsCfg)

	// Injeção de Dependências (Wiring)
	userRepo := dynamo.NewUserRepository(dynamoClient, cfg.DynamoTableName)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	// Setupo do Router e Handlers
	r := api.SetupRouter(authHandler)

	// Start
	slog.Info("Server starting", "port", cfg.Port, "env", cfg.DynamoTableName)
	if err := r.Run(":" + cfg.Port); err != nil {
		slog.Error("Server failed to start", "error", err)
	}
}