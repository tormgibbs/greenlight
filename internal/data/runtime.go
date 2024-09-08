package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Runtime type for runtime of a movie
type Runtime int32

// ErrInvalidRuntimeFormat An error which is returned if UnmarshalJSON() can't parse JSON string successfully
var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

// MarshalJSON method for the Runtime type
func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)

	// Surrounds the string with double quotes
	quotedJSONValue := strconv.Quote(jsonValue)

	// Convert string to value to a byte slice and return it
	return []byte(quotedJSONValue), nil
}

func (r *Runtime) UnmarshalJSON(data []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(data))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(unquotedJSONValue, " ")

	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(i)

	return nil
}
