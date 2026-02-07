package service

import (
	"errors"
	"hack-auth/internal/domain"
	"hack-auth/internal/tests/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestSignUp(t *testing.T) {
	t.Run("should create user successfully", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := NewAuthService(repo, "secret")

		repo.On("FindByEmail", "new@test.com").Return(nil, errors.New("not found"))
		repo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

		err := svc.SignUp("Test", "new@test.com", "123456")

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("should fail if email exists", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := NewAuthService(repo, "secret")

		existingUser := &domain.User{Email: "exists@test.com"}
		repo.On("FindByEmail", "exists@test.com").Return(existingUser, nil)

		err := svc.SignUp("Test", "exists@test.com", "123456")

		assert.Error(t, err)
		assert.Equal(t, "email already registered", err.Error())
	})

	// Teste de erro no repositório ao criar
	t.Run("should fail if repository fails to create", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := NewAuthService(repo, "secret")

		repo.On("FindByEmail", "error@test.com").Return(nil, errors.New("not found"))
		// Simula erro ao salvar no banco
		repo.On("Create", mock.AnythingOfType("*domain.User")).Return(errors.New("dynamo error"))

		err := svc.SignUp("Test", "error@test.com", "123456")

		assert.Error(t, err)
		assert.Equal(t, "dynamo error", err.Error())
	})
}

func TestLogin(t *testing.T) {
	t.Run("should return token on success", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := NewAuthService(repo, "my-secret-key")

		password := "123456"
		hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		hashedPass := string(hashedBytes)

		user := &domain.User{
			Email:    "valid@test.com",
			Password: hashedPass,
			Name:     "Test User",
		}

		repo.On("FindByEmail", "valid@test.com").Return(user, nil)

		token, err := svc.Login("valid@test.com", password)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		repo.AssertExpectations(t)
	})

	t.Run("should fail with wrong password", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := NewAuthService(repo, "my-secret-key")

		hashedBytes, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		user := &domain.User{Email: "valid@test.com", Password: string(hashedBytes)}

		repo.On("FindByEmail", "valid@test.com").Return(user, nil)

		token, err := svc.Login("valid@test.com", "wrong")

		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
		assert.Empty(t, token)
	})

	// Teste de usuário inexistente
	t.Run("should fail if user not found", func(t *testing.T) {
		repo := new(mocks.MockUserRepository)
		svc := NewAuthService(repo, "my-secret-key")

		// Repo retorna erro de "não encontrado"
		repo.On("FindByEmail", "notfound@test.com").Return(nil, errors.New("user not found"))

		token, err := svc.Login("notfound@test.com", "123456")

		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
		assert.Empty(t, token)
	})
}