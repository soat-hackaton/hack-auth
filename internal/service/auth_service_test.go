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
		mockRepo := new(mocks.MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret")

		mockRepo.On("FindByEmail", "new@test.com").Return(nil, nil)
		mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

		err := service.SignUp("Test", "new@test.com", "StrongPass!1")

		assert.Nil(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error if email exists", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret")

		existingUser := &domain.User{Email: "exists@test.com"}
		mockRepo.On("FindByEmail", "exists@test.com").Return(existingUser, nil)

		err := service.SignUp("Test", "exists@test.com", "StrongPass!1")

		assert.Equal(t, domain.ErrUserAlreadyExists, err)
	})

	t.Run("should return internal server error if repository fails", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret")

		mockRepo.On("FindByEmail", "test@example.com").Return(nil, nil)
		mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(errors.New("db error"))

		err := service.SignUp("Test User", "test@example.com", "StrongPass!1")

		assert.NotNil(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestLogin(t *testing.T) {
	password := "StrongPass!1"
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashedPassword := string(hashedBytes)

	t.Run("should return token on success", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret")

		user := &domain.User{
			Email:    "valid@test.com",
			Password: hashedPassword,
		}

		mockRepo.On("FindByEmail", "valid@test.com").Return(user, nil)

		token, err := service.Login("valid@test.com", password)

		assert.Nil(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("should fail with wrong password", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret")

		user := &domain.User{
			Email:    "valid@test.com",
			Password: hashedPassword,
		}

		mockRepo.On("FindByEmail", "valid@test.com").Return(user, nil)

		token, err := service.Login("valid@test.com", "wrong_password")

		assert.NotNil(t, err)
		assert.Equal(t, domain.ErrInvalidCredentials, err)
		assert.Empty(t, token)
	})
}