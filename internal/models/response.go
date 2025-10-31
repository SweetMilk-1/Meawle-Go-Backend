package models

import "fmt"

// Response представляет стандартный ответ API
type Response struct {
	IsOk   bool        `json:"is_ok"`
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// Success создает успешный ответ
func Success(result interface{}) Response {
	return Response{
		IsOk:   true,
		Result: result,
	}
}

// Error создает ответ с ошибкой
func Error(err string) Response {
	return Response{
		IsOk:  false,
		Error: err,
	}
}

// Errorf создает ответ с ошибкой с форматированием
func Errorf(format string, args ...interface{}) Response {
	return Error(fmt.Sprintf(format, args...))
}
