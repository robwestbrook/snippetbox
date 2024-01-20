// Package validator is the validation package.
package validator

import (
	"strings"
	"unicode/utf8"
)

// Validator type containing a map of
// validation errors for the form fields.
type Validator struct {
	FieldErrors		map[string]string
}


// Valid function returns true if FieldErrors map
// doesn't contain any entries.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// AddFieldError function adds an error message to the
// FieldErrors map, as long as no entry already exists
func (v *Validator) AddFieldError(key, message string) {
	// Initialize the map first if not already initialized
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	// If it does already exist, add key and message
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField function adds an error message to the
// FieldErrors map only if a validation check is not OK
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// NotBlank function returns true if a value is not an
// empty string
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars function true if value contains no more
// than n characters.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedInt function returns true if a value is in a
// list of permitted integers
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}