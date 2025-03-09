package common

import "github.com/go-playground/validator"

var Validator = validator.New()

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

func ParseResultToMap(result any) map[string]any {
	return map[string]any{"data": result}
}

func BuildError(err error) map[string]any {
	errResponse := &ErrorResponse{
		Error: err.Error(),
	}

	return ParseResultToMap(errResponse)
}
