package service

import (
	"errors"
	"hack-auth/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      domain.UserRepository
	jwtSecret []byte // Segredo para assinar o token
}

// NewAuthService recebe a interface do repositório e a chave secreta (que virá de env var)
func NewAuthService(repo domain.UserRepository, secret string) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: []byte(secret),
	}
}

// SignUp: Regra de cadastro
func (s *AuthService) SignUp(name, email, password string) error {
	// 1. Verifica se já existe (opcional, pois o repo já trata erro de duplicidade,
	// mas aqui podemos ser mais explícitos se quisermos)
	_, err := s.repo.FindByEmail(email)
	if err == nil {
		return errors.New("email already registered")
	}

	// 2. Hash da Senha (Segurança: Nunca salvar texto plano)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// 3. Cria a entidade de Domínio
	user := &domain.User{
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	// 4. Salva usando o repositório
	return s.repo.Create(user)
}

// Login: Regra de autenticação
func (s *AuthService) Login(email, password string) (string, error) {
	// 1. Busca o usuário
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		// Por segurança, mensagem genérica para não revelar se o email existe ou não
		return "", errors.New("invalid credentials")
	}

	// 2. Compara a senha fornecida com o Hash do banco
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 3. Gera o Token JWT
	token, err := s.generateJWT(user)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

// generateJWT: Cria o token assinado (Método privado auxiliar)
func (s *AuthService) generateJWT(user *domain.User) (string, error) {
	// Define as Claims (o conteúdo do token)
	claims := jwt.MapClaims{
		"sub":   user.Email, // Subject (quem é o dono do token)
		"name":  user.Name,  // Opcional: Nome para o front exibir
		"iss":   "fiap-x",   // Issuer (quem emitiu)
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Expira em 24h
	}

	// Cria o token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Assina com a chave secreta
	return token.SignedString(s.jwtSecret)
}

// Garante que AuthService implementa domain.AuthUseCase
var _ domain.AuthUseCase = (*AuthService)(nil)