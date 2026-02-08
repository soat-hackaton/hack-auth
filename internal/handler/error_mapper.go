func mapDomainError(err error) *rest_err.RestErr {
    switch {
    case errors.Is(err, domain.ErrUserAlreadyExists):
        return rest_err.NewBadRequestError("Email already registered")
    
    case errors.Is(err, domain.ErrPasswordTooWeak):
        return rest_err.NewBadRequestError("Password must be at least 9 chars and contain upper, lower, number, special")
        
    case errors.Is(err, domain.ErrInvalidCredentials):
        return rest_err.NewUnauthorizedError("Invalid credentials")
    
	default:
		return rest_err.NewInternalServerError("Internal server error")
    }
}