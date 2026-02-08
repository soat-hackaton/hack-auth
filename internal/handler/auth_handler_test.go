package handler

import (
	"bytes"
	"errors"
	"hack-auth/internal/tests/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSignUpHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return 201 on success", func(t *testing.T) {
		mockSvc := new(mocks.MockAuthService)
		h := NewAuthHandler(mockSvc)

		mockSvc.On("SignUp", "John", "john@test.com", "StrongPassword1@@").Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		reqBody := []byte(`{"name":"John", "email":"john@test.com", "password":"StrongPassword1@@"}`)
		c.Request, _ = http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		h.SignUp(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("should return 400 on invalid input", func(t *testing.T) {
		mockSvc := new(mocks.MockAuthService)
		h := NewAuthHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		reqBody := []byte(`{"name":"John", "email":"john@test.com"}`) // Sem senha
		c.Request, _ = http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(reqBody))

		h.SignUp(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Teste de Conflito (Email já existe)
	t.Run("should return 409 if email already exists", func(t *testing.T) {
		mockSvc := new(mocks.MockAuthService)
		h := NewAuthHandler(mockSvc)

		// Simula erro específico de negócio
		mockSvc.On("SignUp", "John", "exists@test.com", "StrongPassword1@@").Return(errors.New("email already registered"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		reqBody := []byte(`{"name":"John", "email":"exists@test.com", "password":"StrongPassword1@@"}`)
		c.Request, _ = http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		h.SignUp(c)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("should return 500 on internal error", func(t *testing.T) {
		mockSvc := new(mocks.MockAuthService)
		h := NewAuthHandler(mockSvc)

		// Simula erro genérico
		mockSvc.On("SignUp", "John", "error@test.com", "StrongPassword1@@").Return(errors.New("database connection failed"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		reqBody := []byte(`{"name":"John", "email":"error@test.com", "password":"StrongPassword1@@"}`)
		c.Request, _ = http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		h.SignUp(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return 200 and token", func(t *testing.T) {
		mockSvc := new(mocks.MockAuthService)
		h := NewAuthHandler(mockSvc)

		mockSvc.On("Login", "john@test.com", "123456").Return("mocked-token-jwt", nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		reqBody := []byte(`{"email":"john@test.com", "password":"123456"}`)
		c.Request, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Login(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "mocked-token-jwt")
	})

	t.Run("should return 401 on failure", func(t *testing.T) {
		mockSvc := new(mocks.MockAuthService)
		h := NewAuthHandler(mockSvc)

		mockSvc.On("Login", "john@test.com", "wrong").Return("", errors.New("invalid"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		reqBody := []byte(`{"email":"john@test.com", "password":"wrong"}`)
		c.Request, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))

		h.Login(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should return 400 on invalid login input", func(t *testing.T) {
		mockSvc := new(mocks.MockAuthService)
		h := NewAuthHandler(mockSvc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		reqBody := []byte(`{"email":"john@test.com"}`) // Sem senha
		c.Request, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))

		h.Login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}