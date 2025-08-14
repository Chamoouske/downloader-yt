package dependencyinjections

import (
	"testing"
)

func TestNewVideoDatabase(t *testing.T) {
	// This test primarily checks if the function returns a non-nil database instance.
	// The actual behavior of the database (Save, Get, Remove) is tested in memoria_test.go
	db := NewVideoDatabase()

	if db == nil {
		t.Error("expected a database instance, got nil")
	}
}