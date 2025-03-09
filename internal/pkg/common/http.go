package common

import "github.com/go-playground/validator"

var Validator = validator.New()

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

func BuildError(err error) *ErrorResponse {
	return &ErrorResponse{
		Error: err.Error(),
	}
}
