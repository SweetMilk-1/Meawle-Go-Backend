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

// CatBreedHandler представляет хэндлер для работы с породами кошек
type CatBreedHandler struct {
	service *services.CatBreedService
}

// NewCatBreedHandler создает новый экземпляр хэндлера пород кошек
func NewCatBreedHandler(service *services.CatBreedService) *CatBreedHandler {
	return &CatBreedHandler{service: service}
}

// Create обрабатывает создание породы кошек
func (h *CatBreedHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	var req models.CatBreedCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rw.Error(http.StatusBadRequest, "Invalid request body")
		return
	}

	breed, err := h.service.Create(&req, currentUser.UserID)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	rw.Created(breed)
}

// GetCatBreed обрабатывает получение породы кошек по ID
func (h *CatBreedHandler) GetCatBreed(w http.ResponseWriter, r *http.Request) {
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
		rw.Error(http.StatusBadRequest, "Invalid cat breed ID")
		return
	}

	breed, err := h.service.GetCatBreedByID(id)
	if err != nil {
		rw.Error(http.StatusNotFound, "Cat breed not found")
		return
	}

	rw.Success(breed)
}

// GetAllCatBreeds обрабатывает получение всех пород кошек
func (h *CatBreedHandler) GetAllCatBreeds(w http.ResponseWriter, r *http.Request) {
	rw := NewResponseWriter(w)

	if !ValidateMethod(r, http.MethodGet) {
		rw.Error(ErrMethodNotAllowed.StatusCode, ErrMethodNotAllowed.Message)
		return
	}

	breeds, err := h.service.GetAllCatBreeds()
	if err != nil {
		rw.Error(http.StatusInternalServerError, "Internal server error")
		return
	}

	rw.Success(breeds)
}

// UpdateCatBreed обрабатывает обновление породы кошек
func (h *CatBreedHandler) UpdateCatBreed(w http.ResponseWriter, r *http.Request) {
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
		rw.Error(http.StatusBadRequest, "Invalid cat breed ID")
		return
	}

	var req models.CatBreedUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rw.Error(http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.service.UpdateCatBreed(id, &req, currentUser.UserID, currentUser.IsAdmin)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	rw.Success("Cat breed updated successfully")
}

// DeleteCatBreed обрабатывает удаление породы кошек
func (h *CatBreedHandler) DeleteCatBreed(w http.ResponseWriter, r *http.Request) {
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
		rw.Error(http.StatusBadRequest, "Invalid cat breed ID")
		return
	}

	err = h.service.DeleteCatBreed(id, currentUser.UserID, currentUser.IsAdmin)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	rw.Success("Cat breed deleted successfully")
}

// handleServiceError обрабатывает ошибки сервиса
func (h *CatBreedHandler) handleServiceError(rw *ResponseWriter, err error) {
	switch err {
	case services.ErrCatBreedNotFound:
		rw.Error(http.StatusNotFound, "Cat breed not found")
	case services.ErrCatBreedNameExists:
		rw.Error(http.StatusConflict, "Cat breed name already exists")
	case services.ErrInvalidCatBreedData:
		rw.Error(http.StatusBadRequest, "Invalid cat breed data")
	case services.ErrInvalidCreationDate:
		rw.Error(http.StatusBadRequest, "Invalid creation date")
	case services.ErrAccessDenied:
		rw.Error(http.StatusForbidden, "Access denied")
	default:
		rw.Error(http.StatusInternalServerError, "Internal server error")
	}
}
