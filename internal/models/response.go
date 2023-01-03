package models

type ResponseModel struct {
	Message   string `json:"message"`
	Error     bool   `json:"error"`
	ErrorCode int    `json:"error_code"`
	Data      any    `json:"data,omitempty"`
}
