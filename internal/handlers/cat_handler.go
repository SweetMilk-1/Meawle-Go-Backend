package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"meawle/internal/middleware"
	"meawle/internal/models"
	"meawle/internal/services"

	"github.com/gorilla/mux"
)

// CatHandler представляет хэндлер для работы с котами
type CatHandler struct {
	service *services.CatService
}

// NewCatHandler создает новый экземпляр хэндлера котов
func NewCatHandler(service *services.CatService) *CatHandler {
	return &CatHandler{service: service}
}

// Create обрабатывает создание кота
func (h *CatHandler) Create(w http.ResponseWriter, r *http.Request) {
	rw := NewResponseWriter(w)

	if !ValidateMethod(r, http.MethodPost) {
		rw.Error(ErrMethodNotAllowed.StatusCode, ErrMethodNotAllowed.Message)
		return
	}

	// Получаем пользователя из контекста
	currentUser := middleware.GetUserFromContext(r.Context())
	if currentUser == nil {
		rw.Error(http.StatusUnauthorized, "Authentication required")
		return
	}

	var req models.CatCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rw.Error(http.StatusBadRequest, "Invalid request body")
		return
	}

	cat, err := h.service.Create(&req, currentUser.UserID)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	rw.Created(cat)
}

// GetCat обрабатывает получение кота по ID
func (h *CatHandler) GetCat(w http.ResponseWriter, r *http.Request) {
	rw := NewResponseWriter(w)

	if !ValidateMethod(r, http.MethodGet) {
		rw.Error(ErrMethodNotAllowed.StatusCode, ErrMethodNotAllowed.Message)
		return
	}

	// Извлекаем ID из path параметров
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		rw.Error(http.StatusBadRequest, "Invalid cat ID")
		return
	}

	cat, err := h.service.GetCatByID(id)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	rw.Success(cat)
}

// GetAllCats обрабатывает получение всех котов
func (h *CatHandler) GetAllCats(w http.ResponseWriter, r *http.Request) {
	rw := NewResponseWriter(w)

	if !ValidateMethod(r, http.MethodGet) {
		rw.Error(ErrMethodNotAllowed.StatusCode, ErrMethodNotAllowed.Message)
		return
	}

	cats, err := h.service.GetAllCats()
	if err != nil {
		rw.Error(http.StatusInternalServerError, "Internal server error")
		return
	}

	rw.Success(cats)
}

// GetUserCats обрабатывает получение котов текущего пользователя
func (h *CatHandler) GetUserCats(w http.ResponseWriter, r *http.Request) {
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

	cats, err := h.service.GetUserCats(currentUser.UserID)
	if err != nil {
		rw.Error(http.StatusInternalServerError, "Internal server error")
		return
	}

	rw.Success(cats)
}

// UpdateCat обрабатывает обновление кота
func (h *CatHandler) UpdateCat(w http.ResponseWriter, r *http.Request) {
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

	// Извлекаем ID из path параметров
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		rw.Error(http.StatusBadRequest, "Invalid cat ID")
		return
	}

	var req models.CatUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rw.Error(http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.service.UpdateCat(id, &req, currentUser.UserID, currentUser.IsAdmin)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	rw.Success("Cat updated successfully")
}

// DeleteCat обрабатывает удаление кота
func (h *CatHandler) DeleteCat(w http.ResponseWriter, r *http.Request) {
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

	// Извлекаем ID из path параметров
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		rw.Error(http.StatusBadRequest, "Invalid cat ID")
		return
	}

	err = h.service.DeleteCat(id, currentUser.UserID, currentUser.IsAdmin)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	rw.Success("Cat deleted successfully")
}

// handleServiceError обрабатывает ошибки сервиса
func (h *CatHandler) handleServiceError(rw *ResponseWriter, err error) {
	switch err {
	case services.ErrCatNotFound:
		rw.Error(http.StatusNotFound, "Cat not found")
	case services.ErrInvalidCatData:
		rw.Error(http.StatusBadRequest, "Invalid cat data")
	case services.ErrInvalidCatAge:
		rw.Error(http.StatusBadRequest, "Cat age must be between 0 and 30 years")
	case services.ErrAccessDenied:
		rw.Error(http.StatusForbidden, "Access denied")
	default:
		rw.Error(http.StatusInternalServerError, "Internal server error")
	}
}
