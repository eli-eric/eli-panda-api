package testsetup

import (
	"os"
	"testing"
)

// TestMain is the global test entry point
func TestMain(m *testing.M) {
	// Initialize shared test database
	InitTestDatabase()

	// Run all tests
	code := m.Run()

	// Cleanup after all tests
	CloseTestDatabase()
	os.Exit(code)
}
