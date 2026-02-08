package rest_err

import "net/http"

// RestErr define o contrato padrão de erros da sua API
type RestErr struct {
	Message string `json:"message"`
	Err     string `json:"error"`
	Code    int    `json:"code"`
	Causes  []any  `json:"causes,omitempty"` // Para erros de validação com múltiplos campos
}

// Implementa a interface 'error' do Go
func (r *RestErr) Error() string {
	return r.Message
}

// NewRestErr cria um erro genérico
func NewRestErr(message, err string, code int, causes []any) *RestErr {
	return &RestErr{
		Message: message,
		Err:     err,
		Code:    code,
		Causes:  causes,
	}
}

// Construtores auxiliares para os erros mais comuns

func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "bad_request",
		Code:    http.StatusBadRequest,
	}
}

func NewBadRequestValidationError(message string, causes []any) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "bad_request",
		Code:    http.StatusBadRequest,
		Causes:  causes,
	}
}

func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "internal_server_error",
		Code:    http.StatusInternalServerError,
	}
}

func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "not_found",
		Code:    http.StatusNotFound,
	}
}

func NewForbiddenError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "forbidden",
		Code:    http.StatusForbidden,
	}
}

func NewUnauthorizedError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "unauthorized",
		Code:    http.StatusUnauthorized,
	}
}