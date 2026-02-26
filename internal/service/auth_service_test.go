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

	t.Run("should return error if password is weak", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret")

		err := service.SignUp("Test", "weak@test.com", "123")

		assert.Equal(t, domain.ErrPasswordTooWeak, err)
		mockRepo.AssertNotCalled(t, "FindByEmail")
	})

	t.Run("should return error if email exists", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret")

		existingUser := &domain.User{Email: "exists@test.com"}
		mockRepo.On("FindByEmail", "exists@test.com").Return(existingUser, nil)

		err := service.SignUp("Test", "exists@test.com", "StrongPass!1")

		assert.Equal(t, domain.ErrUserAlreadyExists, err)
	})

	t.Run("should return error if find by email fails", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret")

		mockRepo.On("FindByEmail", "error@test.com").Return(nil, errors.New("connection refused"))

		err := service.SignUp("Test", "error@test.com", "StrongPass!1")

		assert.Error(t, err)
		assert.Equal(t, "connection refused", err.Error())
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
	testPass := "StrongPass!1"
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(testPass), bcrypt.DefaultCost)
	hashedPassword := string(hashedBytes)

	t.Run("should return token on success", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret")

		user := &domain.User{
			Email:    "valid@test.com",
			Password: hashedPassword,
		}

		mockRepo.On("FindByEmail", "valid@test.com").Return(user, nil)

		token, err := service.Login("valid@test.com", testPass)

		assert.Nil(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("should return error if repository fails", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret")

		mockRepo.On("FindByEmail", "error@test.com").Return(nil, errors.New("db error"))

		token, err := service.Login("error@test.com", testPass)

		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("should fail if user not found", func(t *testing.T) {
		mockRepo := new(mocks.MockUserRepository)
		service := NewAuthService(mockRepo, "test-secret")

		mockRepo.On("FindByEmail", "notfound@test.com").Return(nil, nil)

		token, err := service.Login("notfound@test.com", "anypass")

		assert.Equal(t, domain.ErrInvalidCredentials, err)
		assert.Empty(t, token)
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
