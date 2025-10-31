package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"meawle/internal/models"
)

// ResponseWriter представляет обертку для стандартного ResponseWriter
type ResponseWriter struct {
	w http.ResponseWriter
}

// NewResponseWriter создает новый ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w: w}
}

// JSON устанавливает заголовок Content-Type и кодирует данные в JSON
func (rw *ResponseWriter) JSON(statusCode int, data interface{}) {
	rw.w.Header().Set("Content-Type", "application/json")
	rw.w.WriteHeader(statusCode)
	json.NewEncoder(rw.w).Encode(data)
}

// Error отправляет JSON ошибку
func (rw *ResponseWriter) Error(statusCode int, message string) {
	rw.JSON(statusCode, models.Error(message))
}

// Success отправляет JSON успешного ответа
func (rw *ResponseWriter) Success(data interface{}) {
	rw.JSON(http.StatusOK, models.Success(data))
}

// Created отправляет JSON ответа с кодом 201
func (rw *ResponseWriter) Created(data interface{}) {
	rw.JSON(http.StatusCreated, models.Success(data))
}

// ParseID извлекает и парсит ID из query параметров
func ParseID(r *http.Request, paramName string) (int, error) {
	idStr := r.URL.Query().Get(paramName)
	if idStr == "" {
		return 0, ErrMissingParameter
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, ErrInvalidParameter
	}
	
	return id, nil
}

// ValidateMethod проверяет HTTP метод
func ValidateMethod(r *http.Request, allowedMethod string) bool {
	return r.Method == allowedMethod
}

// Ошибки
var (
	ErrMissingParameter = &HandlerError{Message: "Parameter is required", StatusCode: http.StatusBadRequest}
	ErrInvalidParameter = &HandlerError{Message: "Invalid parameter", StatusCode: http.StatusBadRequest}
	ErrMethodNotAllowed = &HandlerError{Message: "Method not allowed", StatusCode: http.StatusMethodNotAllowed}
)

// HandlerError представляет ошибку хэндлера
type HandlerError struct {
	Message    string
	StatusCode int
}

func (e *HandlerError) Error() string {
	return e.Message
}