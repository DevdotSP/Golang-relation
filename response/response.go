package response

// ErrorModel is the structure for API error responses
type ErrorModel struct {
	RetCode any `json:"ret_code"` // Return Code
	Message any `json:"message"`  // Error Message
	Data    any `json:"data"`     // Error details
}
