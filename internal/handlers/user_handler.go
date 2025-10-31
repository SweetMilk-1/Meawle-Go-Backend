package handlers

import (
	"encoding/json"
	"net/http"

	"meawle/internal/middleware"
	"meawle/internal/models"
	"meawle/internal/services"
)

// UserHandler представляет хэндлер для работы с пользователями
type UserHandler struct {
	service *services.UserService
}

// NewUserHandler создает новый экземпляр хэндлера пользователей
func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Register обрабатывает регистрацию пользователя
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	rw := NewResponseWriter(w)

	if !ValidateMethod(r, http.MethodPost) {
		rw.Error(ErrMethodNotAllowed.StatusCode, ErrMethodNotAllowed.Message)
		return
	}

	var req models.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rw.Error(http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.service.Register(&req)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	rw.Created(user)
}

// Login обрабатывает вход пользователя
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	rw := NewResponseWriter(w)

	if !ValidateMethod(r, http.MethodPost) {
		rw.Error(ErrMethodNotAllowed.StatusCode, ErrMethodNotAllowed.Message)
		return
	}

	var req models.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rw.Error(http.StatusBadRequest, "Invalid request body")
		return
	}

	token, user, err := h.service.Login(&req)
	if err != nil {
		rw.Error(http.StatusUnauthorized, "Invalid credentials")
		return
	}

	response := map[string]interface{}{
		"token": token,
		"user":  user,
	}

	rw.Success(response)
}

// GetUser обрабатывает получение пользователя по ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	rw := NewResponseWriter(w)

	if !ValidateMethod(r, http.MethodGet) {
		rw.Error(ErrMethodNotAllowed.StatusCode, ErrMethodNotAllowed.Message)
		return
	}

	// Получаем пользователя из контекста
	currentUser := middleware.GetUserFromContext(r.Context())
	if currentUser == nil {
		rw.Error(http.StatusUnauthorized, "Authentication required")
		return
	}

	// Извлекаем ID из URL параметров
	id, err := ParseID(r, "id")
	if err != nil {
		rw.Error(http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.GetUserByID(id, currentUser.UserID, currentUser.IsAdmin)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	rw.Success(user)
}

// GetAllUsers обрабатывает получение всех пользователей
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	rw := NewResponseWriter(w)

	if !ValidateMethod(r, http.MethodGet) {
		rw.Error(ErrMethodNotAllowed.StatusCode, ErrMethodNotAllowed.Message)
		return
	}

	users, err := h.service.GetAllUsers()
	if err != nil {
		rw.Error(http.StatusInternalServerError, "Internal server error")
		return
	}

	rw.Success(users)
}

// UpdateUser обрабатывает обновление пользователя
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	rw := NewResponseWriter(w)

	if !ValidateMethod(r, http.MethodPut) {
		rw.Error(ErrMethodNotAllowed.StatusCode, ErrMethodNotAllowed.Message)
		return
	}

	// Получаем пользователя из контекста
	currentUser := middleware.GetUserFromContext(r.Context())
	if currentUser == nil {
		rw.Error(http.StatusUnauthorized, "Authentication required")
		return
	}

	// Извлекаем ID из URL параметров
	id, err := ParseID(r, "id")
	if err != nil {
		rw.Error(http.StatusBadRequest, err.Error())
		return
	}

	var req models.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rw.Error(http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.service.UpdateUser(id, &req, currentUser.UserID, currentUser.IsAdmin)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	rw.Success("User updated successfully")
}

// DeleteUser обрабатывает удаление пользователя
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	rw := NewResponseWriter(w)

	if !ValidateMethod(r, http.MethodDelete) {
		rw.Error(ErrMethodNotAllowed.StatusCode, ErrMethodNotAllowed.Message)
		return
	}

	// Получаем пользователя из контекста
	currentUser := middleware.GetUserFromContext(r.Context())
	if currentUser == nil {
		rw.Error(http.StatusUnauthorized, "Authentication required")
		return
	}

	// Извлекаем ID из URL параметров
	id, err := ParseID(r, "id")
	if err != nil {
		rw.Error(http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.DeleteUser(id, currentUser.UserID, currentUser.IsAdmin)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	rw.Success("User deleted successfully")
}

// handleServiceError обрабатывает ошибки сервиса
func (h *UserHandler) handleServiceError(rw *ResponseWriter, err error) {
	switch err {
	case services.ErrUserNotFound:
		rw.Error(http.StatusNotFound, "User not found")
	case services.ErrEmailExists:
		rw.Error(http.StatusConflict, "Email already exists")
	case services.ErrAccessDenied:
		rw.Error(http.StatusForbidden, "Access denied")
	default:
		rw.Error(http.StatusInternalServerError, "Internal server error")
	}
}
