package provider

import "strings"

// isNotFoundError checks if the error indicates a resource was not found.
// The Resend SDK returns generic error strings, so we check the message.
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "not found") || strings.Contains(msg, "not_found")
}
