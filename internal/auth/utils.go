package auth

import "strings"

// ?: does it make sense to make a util function here? or should I just put it in my util?
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
