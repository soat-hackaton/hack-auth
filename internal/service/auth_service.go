package service

import (
	"hack-auth/internal/config"
	"hack-auth/internal/domain"
	"hack-auth/internal/repository/dynamo"
	"hack-auth/internal/utils/rest_err" // Import novo
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const TokenDuration = time.Hour * 24

type authService struct {
	repo dynamo.AuthRepository
}

func NewAuthService(repo dynamo.AuthRepository) domain.AuthUseCase {
	return &authService{repo: repo}
}

func (s *authService) SignUp(name, email, password string) error {
    // 1. Validação da complexidade da senha
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

	// 4. Cria usuário
    return s.repo.CreateUser(user)
}

func (s *authService) Login(email, password string) (string, *rest_err.RestErr) {
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
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(TokenDuration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.Envs.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}