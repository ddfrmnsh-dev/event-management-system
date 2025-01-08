package modelutil

// APIResponse represents a standard API response.
// @Description Standard API response format
type Response struct {
	Message string      `json:"message"`
	Status  bool        `json:"status"`
	Data    interface{} `json:"data"`
}

func APIResponse(message string, data interface{}, status bool) Response {
	jsonResponse := Response{
		Message: message,
		Status:  status,
		Data:    data,
	}

	return jsonResponse
}
