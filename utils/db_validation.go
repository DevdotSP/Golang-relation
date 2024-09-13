package utils

import "strings"

// Helper function to detect unique constraint violation
func IsUniqueConstraintError(err error) bool {
	// Check if error contains specific keywords indicating a unique constraint violation
	return err != nil && (strings.Contains(err.Error(), "unique constraint"))
}