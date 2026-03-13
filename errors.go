package normify

import "fmt"

type APIError struct {
	StatusCode int
	Status     string
	Body       []byte
}

func (e *APIError) Error() string {
	return fmt.Sprintf("normify API error: %s (status %d)", e.Status, e.StatusCode)
}
