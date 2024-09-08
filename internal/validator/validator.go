package validator

import (
	"regexp"
)

var (
	EmailRX = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+\/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
)

// Validator holds validation errors. It uses a map where the key is a field name
// and the value is the error message associated with that field.
type Validator struct {
	Errors map[string]string
}

// New creates and returns a new Validator instance with an initialized Errors map.
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if there are no errors in the Errors map, indicating that all validations passed.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message to the Errors map for a given field if the error doesn't already exist.
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error message to the Errors map if the given condition `ok` is false.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// In checks if a given value is present in a list of strings. It returns true if the value is found.
func In(value string, list ...string) bool {
	for _, v := range list {
		if value == v {
			return true
		}
	}
	return false
}

// Matches checks if a given value matches a regular expression. It returns true if the value matches.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique checks if all values in a slice are unique. It returns true if all values are unique.
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
