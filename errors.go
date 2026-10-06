package kkhay

import "fmt"

// APIError represents an error response returned by the K Khay API.
type APIError struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Code    string      `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("kkhay: API error HTTP %d [%s]: %s", e.Status, e.Code, e.Message)
	}
	return fmt.Sprintf("kkhay: API error HTTP %d: %s", e.Status, e.Message)
}

