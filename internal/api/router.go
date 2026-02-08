package api

import (
	"hack-auth/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter configura todas as rotas da aplicação
func SetupRouter(authHandler *handler.AuthHandler) *gin.Engine {
	r := gin.Default()

	// Middlewares globais (CORS, Logger, etc) podem vir aqui
	// r.Use(middleware.Cors())

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Routes
	r.POST("/signup", authHandler.SignUp)
	r.POST("/login", authHandler.Login)

	return r
}