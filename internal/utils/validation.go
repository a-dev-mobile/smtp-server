// internal/utils/validation.go

package utils

import (
    "regexp"
)

// ValidateEmail проверяет, является ли строка допустимым email адресом.
func ValidateEmail(email string) bool {
    // Строгая регулярная проверка для email может быть сложной, но этот пример достаточно хорош для начала.
    re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    return re.MatchString(email)
}
