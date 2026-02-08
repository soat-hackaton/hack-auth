package logger

import (
	"log/slog"
	"os"
)

// Init configura o logger padrão para escrever JSON no Stdout
func Init() {
	// Cria um handler JSON.
	// AddSource: true adiciona o arquivo e linha onde o log ocorreu (ótimo para debug)
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true, 
	})

	// Define como logger global
	logger := slog.New(handler)
	slog.SetDefault(logger)
}