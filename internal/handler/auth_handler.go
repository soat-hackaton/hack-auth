package handler

import (
	"fmt"
	"hack-auth/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// DTOs (Data Transfer Objects) para validação de entrada
type signUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"` // aumentar complexidade
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// SignUp - Rota de Cadastro
func (h *AuthHandler) SignUp(c *gin.Context) {
	var req signUpRequest

	// 1. BindJSON faz o parse e valida as tags (required, email, min)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Chama o serviço
	err := h.authService.SignUp(req.Name, req.Email, req.Password)
	if err != nil {
		fmt.Printf("ERRO NO SIGNUP: %v\n", err)

		// Se o erro for "email já existe", retornamos 409 Conflict
		if err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}
		
		// Outros erros (banco fora do ar, etc)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// Login - Rota de Autenticação
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		// Por segurança, sempre 401 Unauthorized para erro de login
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   86400, // 24h em segundos
	})
}