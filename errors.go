package smail

import "fmt"

// APIError is returned when the Smail API responds with a non-2xx status code.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("smail: HTTP %d: %s", e.StatusCode, e.Message)
}
