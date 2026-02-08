package validator

import "regexp"

// IsPasswordStrong verifica se a senha atende aos requisitos de complexidade:
// - Pelo menos 8 caracteres
// - 1 Letra Maiúscula
// - 1 Letra Minúscula
// - 1 Número
// - 1 Caractere Especial
func IsPasswordStrong(password string) bool {
	var (
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[!@#~$%^&*()+|_]`).MatchString(password)
	)
	
	return hasUpper && hasLower && hasNumber && hasSpecial
}