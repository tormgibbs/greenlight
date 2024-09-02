package data

import (
	"fmt"
	"strconv"
)

// A custom Runtime type
type Runtime int32

// A MarshalJSON() method for the Runtime type
func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)

	// Surrounds the string with double quotes
	quotedJSONValue := strconv.Quote(jsonValue)

	// Convert string to value to a byte slice and return it
	return []byte(quotedJSONValue), nil
}
