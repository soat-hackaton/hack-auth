package domain

type AuthUseCase interface {
	SignUp(name, email, password string) error
	Login(email, password string) (string, error)
}