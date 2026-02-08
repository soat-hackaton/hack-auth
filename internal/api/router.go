package api

import (
	"hack-auth/internal/api/middleware"
	"hack-auth/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter configura todas as rotas da aplicação
func SetupRouter(authHandler *handler.AuthHandler) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.StructuredLogger())

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Routes
	r.POST("/signup", authHandler.SignUp)
	r.POST("/login", authHandler.Login)

	return r
}