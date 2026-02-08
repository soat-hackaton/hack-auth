package handler

import (
	"fmt"
	"hack-auth/internal/domain"
	"hack-auth/internal/utils/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService domain.AuthUseCase
}

func NewAuthHandler(authService domain.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// SignUp - Rota de Cadastro
func (h *AuthHandler) SignUp(c *gin.Context) {
	var req signUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !validator.IsPasswordStrong(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must contain uppercase, lowercase, number and special char",
		})
		return
	}

	err := h.authService.SignUp(req.Name, req.Email, req.Password)
	if err != nil {
		if err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, authResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   86400,
	})
}