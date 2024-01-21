// Package validator is the validation package.
package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Validator type.
// Contains:
//	1. NonFieldErrors - validation errors slice not
//											related to specific form fields
//	2. FieldErrors - 	validation errors map related to
//										specific form fields
type Validator struct {
	NonFieldErrors	[]string
	FieldErrors			map[string]string
}

// Use the regexp.MustCompile() function to parse a
// regular expression pattern for sanity checking of
// an email address. A pointer is returned to a
// "compiled" regexp.Regexp type, or panics on error.
// The "compiled" regexp.Regexp pattern is stored
// in the variable EmailRX.
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zAZ0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")


// Valid function returns true if FieldErrors map
// and NonFieldErrors slice don't contain any entries.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 &&
					len(v.NonFieldErrors) == 0
}

// AddNonFieldError function adds an error message to
// the NonFieldError slice.
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
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

// MinChars function returns true if a value contains
// at least n characters.
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// Matches function returns true if a value matches a
// provided compiled regular expression pattern.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}