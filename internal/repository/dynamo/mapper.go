package dynamo

import (
	"fmt"
	"hack-auth/internal/domain"
)

// toSchema converte do Domain para o Schema do Banco
func toSchema(u *domain.User) UserSchema {
	return UserSchema{
		PK:        fmt.Sprintf("USER#%s", u.Email),
		SK:        "PROFILE",
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
	}
}

// toDomain converte do Schema do Banco para o Domain
func toDomain(schema UserSchema) *domain.User {
	return &domain.User{
		Name:      schema.Name,
		Email:     schema.Email,
		Password:  schema.Password,
		CreatedAt: schema.CreatedAt,
	}
}