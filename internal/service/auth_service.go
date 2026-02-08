package service

import (
	"hack-auth/internal/domain"
	"hack-auth/internal/utils/validator"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const TokenDuration = time.Hour * 24

type authService struct {
	repo      domain.UserRepository
	jwtSecret string
}

func NewAuthService(repo domain.UserRepository, secret string) domain.AuthUseCase {
	return &authService{
		repo:      repo,
		jwtSecret: secret,
	}
}

func (s *authService) SignUp(name, email, password string) error {
	// 1. Validação de Regra de Negócio
	if !validator.IsPasswordStrong(password) {
		return domain.ErrPasswordTooWeak
	}

	// 2. Verifica existência
	existingUser, err := s.repo.FindByEmail(email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return domain.ErrUserAlreadyExists
	}

	// 3. Hash e Persistência
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	return s.repo.Create(user)
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.Email,
		"email": user.Email,
		"exp":   time.Now().Add(TokenDuration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}