package mocks

import (
	"hack-auth/internal/domain"
	"hack-auth/internal/utils/rest_err"

	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) SignUp(name, email, password string) *rest_err.RestErr {
	args := m.Called(name, email, password)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*rest_err.RestErr)
}

func (m *MockAuthService) Login(email, password string) (string, *rest_err.RestErr) {
	args := m.Called(email, password)
	
	token := args.String(0)
	
	if args.Get(1) == nil {
		return token, nil
	}
	return token, args.Get(1).(*rest_err.RestErr)
}

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) CreateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockAuthRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}