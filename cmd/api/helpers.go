package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Define an envelope type
type envelope map[string]interface{}

// Read ID parameter from request
func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// writeJSON() sends a JSON response to the client. It encodes the given data to JSON,
// sets the "Content-Type: application/json" header, writes the provided HTTP status code,
// and adds any additional headers. Returns an error if JSON encoding fails.
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// Encode the data to JSON, returning error if any
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal applications.
	js = append(js, '\n')

	// Set response headers from the provided headers map
	for key, value := range headers {
		w.Header()[key] = value
	}

	// Add the "Content-Type: application/json" header, then write the status code and
	// JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		return err
	}

	return nil
}

// readJSON() reads and decodes JSON from an HTTP request body into the destination struct (dst).
// If there's an issue during decoding, it returns a descriptive error.
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Using http.MaxBytesReader() to limit size of request body to 1MB
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Initializing json.Decoder and calling the DisallowUnknownFields() to return an error if JSON from
	// cannot be mapped to target destination
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Decode the JSON from the request body into the dst variable
	err := dec.Decode(dst)
	if err != nil {
		// Declare variables for specific JSON error types.
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalFieldError *json.InvalidUnmarshalError

		// Use a switch to handle different types of errors that might occur during decoding.
		switch {
		// Check if the error is a syntax error in the JSON and return a descriptive error.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		// Handle the case where the JSON is incomplete or ends unexpectedly.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		// Check if the error is due to a mismatch between the JSON type and the expected Go type.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// Handle empty body scenarios by returning an appropriate error.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// Handling unknown fields
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown field %q", fieldName)
		// Handling a large body request
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		// Panic in case of an internal error with unmarshalling
		case errors.As(err, &invalidUnmarshalFieldError):
			panic(err)

		//	For any other type of error, return the error as is.
		default:
			return err
		}
	}

	// Calls Decode() to check if there is another json value
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
