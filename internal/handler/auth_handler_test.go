package handler

import (
	"bytes"
	"encoding/json"
	"hack-auth/internal/service"
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

		// Mock retorna nil (sucesso)
		mockService.On("SignUp", "Test", "test@example.com", "StrongPassword!1").Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(reqBody))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("should return 400 for weak password", func(t *testing.T) {
		mockService := new(mocks.MockAuthService)
		authHandler := NewAuthHandler(mockService)
		r := gin.Default()
		r.POST("/signup", authHandler.SignUp)

		// Senha fraca (curta ou sem requisitos)
		reqBody := []byte(`{
			"name": "Test",
			"email": "test@example.com",
			"password": "123" 
		}`)

		// O Service NÃO deve ser chamado, pois o validador barra antes
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(reqBody))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response rest_err.RestErr
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "bad_request", response.Err)
        // Valida parte da mensagem
		assert.Contains(t, response.Message, "Password must be at least 9 chars") 
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

		// Mock retorna um RestErr criado pelo service
		errReturn := rest_err.NewBadRequestError("Email already registered")
		mockService.On("SignUp", "Test", "exists@example.com", "StrongPassword!1").Return(errReturn)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(reqBody))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		
		var response rest_err.RestErr
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Email already registered", response.Message)
	})
}

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return token and correct expires_in", func(t *testing.T) {
		mockService := new(mocks.MockAuthService)
		authHandler := NewAuthHandler(mockService)

		r := gin.Default()
		r.POST("/login", authHandler.Login)

		reqBody := []byte(`{
			"email": "test@example.com",
			"password": "StrongPassword!1"
		}`)

		// Mock sucesso
		mockService.On("Login", "test@example.com", "StrongPassword!1").Return("mocked_token", nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "mocked_token", response["access_token"])
		assert.Equal(t, "Bearer", response["token_type"])
		
		// Valida se o expires_in reflete a constante do service
		expectedExpires := int(service.TokenDuration.Seconds())
		assert.Equal(t, float64(expectedExpires), response["expires_in"])
	})

    t.Run("should return 401 on invalid credentials", func(t *testing.T) {
        mockService := new(mocks.MockAuthService)
        authHandler := NewAuthHandler(mockService)
        r := gin.Default()
        r.POST("/login", authHandler.Login)

        reqBody := []byte(`{"email": "test@example.com", "password": "Wrong"}`)

        // Mock retornando erro 401
        errReturn := rest_err.NewUnauthorizedError("Invalid credentials")
        mockService.On("Login", "test@example.com", "Wrong").Return("", errReturn)

        w := httptest.NewRecorder()
        req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
        r.ServeHTTP(w, req)

        assert.Equal(t, http.StatusUnauthorized, w.Code)
        
        var response rest_err.RestErr
        json.Unmarshal(w.Body.Bytes(), &response)
        assert.Equal(t, "unauthorized", response.Err)
    })
}