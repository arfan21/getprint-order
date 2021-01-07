package utils

import (
	"encoding/json"
	"fmt"
)

type ResponseStruct struct {
	Status  string      `json:"status"`
	Message interface{} `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Response(status string, message interface{}, data interface{}) ResponseStruct {
	dataType := fmt.Sprintf("%T", message)

	if dataType == "error" {
		errorJSON, _ := json.Marshal(message)
		return ResponseStruct{status, errorJSON, nil}
	}

	return ResponseStruct{status, message, data}
}
