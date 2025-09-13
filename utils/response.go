package utils

type MessageResponse struct {
	Message string `json:"message"`
}

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
