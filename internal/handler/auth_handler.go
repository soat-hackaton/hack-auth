package handler

import (
	"hack-auth/internal/domain"
	"hack-auth/internal/service"
	"hack-auth/internal/utils/rest_err"
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

func (h *AuthHandler) SignUp(c *gin.Context) {
    var req signUpRequest

    // 1. O Handler converte o JSON (Responsabilidade HTTP)
    if err := c.ShouldBindJSON(&req); err != nil {
        errRest := rest_err.NewBadRequestError("Invalid JSON body")
        c.JSON(errRest.Code, errRest)
        return
    }

    // 2. Chama o Service (Passando DADOS, não contexto HTTP)
    err := h.authService.SignUp(req.Name, req.Email, req.Password)
    if err != nil {
        // Usa nosso Mapper para traduzir o erro de domínio para HTTP
        errRest := mapDomainError(err)
        c.JSON(errRest.Code, errRest)
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errRest := rest_err.NewBadRequestError("Invalid login fields")
		c.JSON(errRest.Code, errRest)
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		errRest := mapDomainError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, authResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(service.TokenDuration.Seconds()),
	})
}