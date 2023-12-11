package server

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/goccy/go-json"
)

// Function for validating an argument (wraps the old CLI validators)
type ArgValidator[ArgType any] func(string, string) (ArgType, error)

// Validates an argument, ensuring it exists and can be converted to the required type
func ValidateArg[ArgType any](name string, args map[string]string, impl ArgValidator[ArgType], result_Out *ArgType) error {
	// Make sure it exists
	arg, exists := args[name]
	if !exists {
		return fmt.Errorf("missing argument '%s'", name)
	}

	// Run the parser
	result, err := impl(name, arg)
	if err != nil {
		return err
	}

	// Set the result
	*result_Out = result
	return nil
}

// Validates an optional argument, converting to the required type if it exists
func ValidateOptionalArg[ArgType any](name string, args map[string]string, impl ArgValidator[ArgType], result_Out *ArgType, exists_Out *bool) error {
	// Make sure it exists
	arg, exists := args[name]
	if !exists {
		if exists_Out != nil {
			*exists_Out = false
		}
		return nil
	}

	// Run the parser
	result, err := impl(name, arg)
	if err != nil {
		return err
	}

	// Set the result
	*result_Out = result
	if exists_Out != nil {
		*exists_Out = true
	}
	return nil
}

// Gets a string argument, ensuring that it exists in the provided vars list
func GetStringFromValues(name string, args url.Values, result_Out *string) error {
	// Make sure it exists
	arg, exists := args[name]
	if !exists {
		return fmt.Errorf("missing argument '%s'", name)
	}

	// Set the result
	*result_Out = arg[0]
	return nil
}

// Gets an optional string argument from the provided vars list
func GetOptionalStringFromVars(name string, args map[string]string, result_Out *string) bool {
	// Make sure it exists
	arg, exists := args[name]
	if !exists {
		return false
	}

	// Set the result
	*result_Out = arg
	return true
}

// Handles an error related to parsing the input parameters of a request
func handleInputError(w http.ResponseWriter, err error) {
	// Write out any errors
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
}

// Handle routes called with an invalid method
func handleInvalidMethod(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte{})
}

// Handles a Node daemon response
func handleResponse(w http.ResponseWriter, response any, err error) {
	// Write out any errors
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	// Write the serialized response
	bytes, err := json.Marshal(response)
	if err != nil {
		err = fmt.Errorf("error serializing response: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	}
}
