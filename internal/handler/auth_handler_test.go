package handler

import (
	"errors"
	"bytes"
	"encoding/json"
	"hack-auth/internal/domain"
	"hack-auth/internal/tests/mocks"
	"hack-auth/internal/utils/rest_err"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSignUpHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return 201 on success", func(t *testing.T) {
		mockService := new(mocks.MockAuthService)
		authHandler := NewAuthHandler(mockService)

		r := gin.Default()
		r.POST("/signup", authHandler.SignUp)

		reqBody := []byte(`{
			"name": "Test",
			"email": "test@example.com",
			"password": "StrongPassword!1" 
		}`)

		mockService.On("SignUp", "Test", "test@example.com", "StrongPassword!1").Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(reqBody))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("should return 400 on invalid json", func(t *testing.T) {
		mockService := new(mocks.MockAuthService)
		authHandler := NewAuthHandler(mockService)
		r := gin.Default()
		r.POST("/signup", authHandler.SignUp)

		// JSON quebrado
		reqBody := []byte(`{"name": "Test", "email": `)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(reqBody))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return 400 if password is weak", func(t *testing.T) {
		mockService := new(mocks.MockAuthService)
		authHandler := NewAuthHandler(mockService)
		r := gin.Default()
		r.POST("/signup", authHandler.SignUp)

		reqBody := []byte(`{"name": "Test", "email": "weak@test.com", "password": "123"}`)
		
		mockService.On("SignUp", "Test", "weak@test.com", "123").Return(domain.ErrPasswordTooWeak)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(reqBody))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "must be at least 8 chars")
	})

	t.Run("should return 400 if email already exists", func(t *testing.T) {
		mockService := new(mocks.MockAuthService)
		authHandler := NewAuthHandler(mockService)
		r := gin.Default()
		r.POST("/signup", authHandler.SignUp)

		reqBody := []byte(`{
			"name": "Test",
			"email": "exists@example.com",
			"password": "StrongPassword!1"
		}`)

		mockService.On("SignUp", "Test", "exists@example.com", "StrongPassword!1").Return(domain.ErrUserAlreadyExists)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(reqBody))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response rest_err.RestErr
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response.Message, "Email already registered")
	})

	t.Run("should return 500 on internal error", func(t *testing.T) {
		mockService := new(mocks.MockAuthService)
		authHandler := NewAuthHandler(mockService)
		r := gin.Default()
		r.POST("/signup", authHandler.SignUp)

		reqBody := []byte(`{"name": "Test", "email": "error@example.com", "password": "StrongPassword!1"}`)
		
		// Mock retorna erro não mapeado (genérico)
		mockService.On("SignUp", "Test", "error@example.com", "StrongPassword!1").Return(errors.New("db dead"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(reqBody))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return token", func(t *testing.T) {
		mockService := new(mocks.MockAuthService)
		authHandler := NewAuthHandler(mockService)
		r := gin.Default()
		r.POST("/login", authHandler.Login)

		reqBody := []byte(`{"email": "test@example.com", "password": "StrongPassword!1"}`)

		mockService.On("Login", "test@example.com", "StrongPassword!1").Return("mocked_token", nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return 400 on invalid json login", func(t *testing.T) {
		mockService := new(mocks.MockAuthService)
		authHandler := NewAuthHandler(mockService)
		r := gin.Default()
		r.POST("/login", authHandler.Login)

		reqBody := []byte(`{"email": "test"`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return 401 on invalid credentials", func(t *testing.T) {
		mockService := new(mocks.MockAuthService)
		authHandler := NewAuthHandler(mockService)
		r := gin.Default()
		r.POST("/login", authHandler.Login)

		reqBody := []byte(`{"email": "test@example.com", "password": "Wrong"}`)

		mockService.On("Login", "test@example.com", "Wrong").Return("", domain.ErrInvalidCredentials)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}