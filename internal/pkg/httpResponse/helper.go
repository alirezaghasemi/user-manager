package httpresponse

// success response
func Success(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// error response with optional trace
func Error(message string, err interface{}) APIResponse {
	return APIResponse{
		Success: false,
		Message: message,
		Errors:  err,
	}
}
