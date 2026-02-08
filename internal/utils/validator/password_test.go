package validator

import "testing"

func TestIsPasswordStrong(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "should return true for strong password",
			password: "StrongPass1!",
			want:     true,
		},
		{
			name:     "should return false if length is less than 8",
			password: "Weak1!",
			want:     false,
		},
		{
			name:     "should return false if missing uppercase",
			password: "weakpass1!",
			want:     false,
		},
		{
			name:     "should return false if missing lowercase",
			password: "WEAKPASS1!",
			want:     false,
		},
		{
			name:     "should return false if missing number",
			password: "NoNumber!",
			want:     false,
		},
		{
			name:     "should return false if missing special character",
			password: "NoSpecial1",
			want:     false,
		},
		{
			name:     "should return true for complex password with symbols",
			password: "@AnoTher$ComPlex99",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPasswordStrong(tt.password); got != tt.want {
				t.Errorf("IsPasswordStrong() = %v, want %v (password: %s)", got, tt.want, tt.password)
			}
		})
	}
}