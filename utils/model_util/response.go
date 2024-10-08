package modelutil

type Response struct {
	Message    string      `json:"message"`
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
}

func APIResponse(message string, data interface{}, statusCode int) Response {
	jsonResponse := Response{
		Message:    message,
		StatusCode: statusCode,
		Data:       data,
	}

	return jsonResponse
}
