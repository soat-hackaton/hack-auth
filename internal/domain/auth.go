package domain

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthUseCase interface {
	SignUp(name, email, password string) error
	Login(email, password string) (string, error)
}