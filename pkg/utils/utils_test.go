package utils

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", "file"},
		{"special chars", "a<b>c:d\\e/f\\g", "a_b_c_d_e_f_g"},
		{"spaces", "a b c ", "a_b_c"},
		{"reserved name", "con", "con_"},
		{"unicode", "ñáñóú", "nanou"},
		{"extension", "file.txt", "file.txt"},
		{"long name", strings.Repeat("a", 300), strings.Repeat("a", 255)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeFilename(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	testKey := "TEST_ENV_VAR"

	os.Unsetenv(testKey)
	assert.Equal(t, "default", GetEnvOrDefault(testKey, "default"))

	os.Setenv(testKey, "custom")
	assert.Equal(t, "custom", GetEnvOrDefault(testKey, "default"))
}
