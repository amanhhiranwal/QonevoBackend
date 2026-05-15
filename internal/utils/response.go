package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func JSON(w http.ResponseWriter, status int, message string, data interface{}, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(Response{
		Success: status >= 200 && status < 300,
		Message: message,
		Data:    data,
		Error:   err,
	})
}