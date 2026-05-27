package response

// ErrorResponse contains error body
type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

// ErrorBody includes error code and message
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewErrorResponse creates a new ErrorResponse
func NewErrorResponse(code string, message string) *ErrorResponse {
	return &ErrorResponse{
		Error: ErrorBody{
			Code:    code,
			Message: message,
		},
	}
}
