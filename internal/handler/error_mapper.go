package handler

import (
	"errors"
	"hack-auth/internal/domain"
	"hack-auth/internal/utils/rest_err"
	"log/slog"
)

func mapDomainError(err error) *rest_err.RestErr {
	switch {
		case errors.Is(err, domain.ErrUserAlreadyExists):
			slog.Warn("Business error", "reason", "email_duplicate", "error", err)
			return rest_err.NewBadRequestError("Email já foi registrado")

		case errors.Is(err, domain.ErrPasswordTooWeak):
			slog.Warn("Business error", "reason", "password_policy", "error", err)
			return rest_err.NewBadRequestError("Senha deve conter no mínimo 8 caracteres, uma maiúscula, uma minúscula, um número e um símbolo")

		case errors.Is(err, domain.ErrInvalidCredentials):
			slog.Warn("Security warning", "reason", "invalid_login", "error", err)
			return rest_err.NewUnauthorizedError("Credenciais inválidas")
		
		case errors.Is(err, domain.ErrUserNotFound):
			slog.Warn("Business error", "reason", "user_not_found", "error", err)
			return rest_err.NewNotFoundError("Usuário não encontrado")

		default:
			slog.Error("Unhandled internal error", "error", err)
			return rest_err.NewInternalServerError("Erro interno no servidor")
	}
}