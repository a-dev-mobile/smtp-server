// internal/utils/validation.go

package utils

import (
    "regexp"
)

// ValidateEmail checks if the string is a valid email address.
func ValidateEmail(email string) bool {
      re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    return re.MatchString(email)
}
