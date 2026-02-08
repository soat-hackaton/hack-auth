package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	JWTSecret       string
	AWSRegion       string
	DynamoTableName string
}

// Load carrega as variáveis de ambiente e valida se as obrigatórias existem
func Load() (*Config, error) {
	// Tenta carregar o arquivo .env (não retorna erro se não existir, 
	// pois em produção/k8s as vars vêm do ambiente direto)
	_ = godotenv.Load()

	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		AWSRegion:       getEnv("AWS_REGION", "us-west-2"),
		DynamoTableName: os.Getenv("DYNAMO_TABLE_NAME"),
	}

	// Validação: Se faltar algo crítico, retornamos erro impedindo o app de subir
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("environment variable JWT_SECRET is required")
	}
	if cfg.DynamoTableName == "" {
		return nil, fmt.Errorf("environment variable DYNAMO_TABLE_NAME is required")
	}

	// Opcional: Validar se credenciais AWS existem (útil no AWS Academy)
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		return nil, fmt.Errorf("AWS credentials (ACCESS_KEY_ID/SECRET_ACCESS_KEY) are missing")
	}

	return cfg, nil
}

// Helper para ler env com valor default
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}