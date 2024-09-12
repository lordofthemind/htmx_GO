package responses

import "time"

// StandardResponse defines the structure for API responses
type StandardResponse struct {
	Status    int         `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     interface{} `json:"error,omitempty"`
	Timestamp string      `json:"timestamp"`
	RequestID string      `json:"requestId,omitempty"`
}

// NewResponse returns a standardized response
func NewResponse(status int, message string, data interface{}, err interface{}) StandardResponse {
	return StandardResponse{
		Status:    status,
		Message:   message,
		Data:      data,
		Error:     err,
		Timestamp: time.Now().Format(time.RFC3339),
		// Optionally, include RequestID from context if needed
	}
}
