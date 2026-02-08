package service

import (
	"errors"
	"hack-auth/internal/domain"
	"hack-auth/internal/tests/mocks"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestSignUp(t *testing.T) {
	t.Run("should create user successfully", func(t *testing.T) {
		mockRepo := new(mocks.MockAuthRepository)
		service := NewAuthService(mockRepo)

		// Mock: FindByEmail retorna nil (usuário não existe)
		mockRepo.On("FindByEmail", "test@example.com").Return(nil, nil)
		// Mock: CreateUser retorna sucesso
		mockRepo.On("CreateUser", mock.Anything).Return(nil)

		// Executa
		err := service.SignUp("Test User", "test@example.com", "Password@123")

		// Valida
		assert.Nil(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return bad request if email exists", func(t *testing.T) {
		mockRepo := new(mocks.MockAuthRepository)
		service := NewAuthService(mockRepo)

		existingUser := &domain.User{Email: "test@example.com"}
		
		// Mock: FindByEmail encontra usuário
		mockRepo.On("FindByEmail", "test@example.com").Return(existingUser, nil)

		// Executa
		err := service.SignUp("Test User", "test@example.com", "Password@123")

		// Valida: Esperamos um RestErr com Code 400
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.Code)
		assert.Equal(t, "Email already registered", err.Message)
		
		// CreateUser NUNCA deve ser chamado
		mockRepo.AssertNotCalled(t, "CreateUser")
	})

	t.Run("should return internal server error if repository fails", func(t *testing.T) {
		mockRepo := new(mocks.MockAuthRepository)
		service := NewAuthService(mockRepo)

		mockRepo.On("FindByEmail", "test@example.com").Return(nil, nil)
		// Mock: CreateUser falha
		mockRepo.On("CreateUser", mock.Anything).Return(errors.New("db error"))

		err := service.SignUp("Test User", "test@example.com", "Password@123")

		// Valida: Esperamos um RestErr com Code 500
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		assert.Equal(t, "Error creating user database record", err.Message)
	})
}

func TestLogin(t *testing.T) {
	// Preparação comum (Setup)
	// Precisamos simular um hash real, pois o Service usa bcrypt.CompareHashAndPassword
	password := "123456"
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedPassword := string(hashedBytes)

	t.Run("should return token on success", func(t *testing.T) {
		mockRepo := new(mocks.MockAuthRepository)
		// Nota: NewAuthService não recebe a key, ela vem do config interno
		service := NewAuthService(mockRepo) 

		user := &domain.User{
			ID:       "user-id-123",
			Email:    "valid@test.com",
			Password: hashedPassword, // O banco retorna a senha criptografada
			Name:     "Test User",
		}

		// Mock: Encontra o usuário
		mockRepo.On("FindByEmail", "valid@test.com").Return(user, nil)

		// Ação
		token, err := service.Login("valid@test.com", password)

		// Validação
		assert.Nil(t, err)          // Erro deve ser nulo (ponteiro nil)
		assert.NotEmpty(t, token)   // Token deve vir preenchido
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail with wrong password", func(t *testing.T) {
		mockRepo := new(mocks.MockAuthRepository)
		service := NewAuthService(mockRepo)

		user := &domain.User{
			Email:    "valid@test.com",
			Password: hashedPassword,
		}

		// Mock: Encontra o usuário, mas a senha enviada no Login será errada
		mockRepo.On("FindByEmail", "valid@test.com").Return(user, nil)

		// Ação: Passamos "wrong_password"
		token, err := service.Login("valid@test.com", "wrong_password")

		// Validação
		assert.NotNil(t, err) // Erro não deve ser nulo
		assert.Empty(t, token)
		
		// Verificamos o Status Code e a Mensagem do RestErr
		assert.Equal(t, http.StatusUnauthorized, err.Code)
		assert.Equal(t, "Invalid credentials", err.Message)
		assert.Equal(t, "unauthorized", err.Err)
	})

	t.Run("should fail if user not found", func(t *testing.T) {
		mockRepo := new(mocks.MockAuthRepository)
		service := NewAuthService(mockRepo)

		// Mock: Retorna nil (usuário não encontrado)
		// Nota: Dependendo da sua implementação do Repo, pode retornar (nil, nil) ou (nil, erro)
		// Assumindo que o repo retorna (nil, nil) quando não acha:
		mockRepo.On("FindByEmail", "notfound@test.com").Return(nil, nil)

		token, err := service.Login("notfound@test.com", "123456")

		// Validação
		assert.NotNil(t, err)
		assert.Empty(t, token)
		
		// Service deve retornar 401 para não expor que o email não existe (segurança)
		assert.Equal(t, http.StatusUnauthorized, err.Code) 
		assert.Equal(t, "Invalid credentials", err.Message)
	})

	t.Run("should return internal server error if db fails", func(t *testing.T) {
		mockRepo := new(mocks.MockAuthRepository)
		service := NewAuthService(mockRepo)

		// Mock: Banco dá erro de conexão
		mockRepo.On("FindByEmail", "error@test.com").Return(nil, errors.New("db connection failed"))

		token, err := service.Login("error@test.com", "123456")

		// Validação
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		assert.Equal(t, "Error validating user", err.Message)
	})
}