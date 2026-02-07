package mocks

import (
	"hack-auth/internal/domain"
	"github.com/stretchr/testify/mock"
)

// --- Mock do Repositório (Para testar o Service) ---
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// --- Mock do Service (Para testar o Handler) ---
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) SignUp(name, email, password string) error {
	args := m.Called(name, email, password)
	return args.Error(0)
}

func (m *MockAuthService) Login(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}