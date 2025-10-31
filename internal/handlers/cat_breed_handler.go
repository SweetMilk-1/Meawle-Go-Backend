package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"meawle/internal/middleware"
	"meawle/internal/models"
	"meawle/internal/services"
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем пользователя из контекста
	currentUser := middleware.GetUserFromContext(r.Context())
	if currentUser == nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req models.CatBreedCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	breed, err := h.service.Create(&req, currentUser.UserID)
	if err != nil {
		switch err {
		case services.ErrCatBreedNameExists:
			http.Error(w, "Cat breed name already exists", http.StatusConflict)
		case services.ErrInvalidCatBreedData:
			http.Error(w, "Invalid cat breed data", http.StatusBadRequest)
		case services.ErrInvalidCreationDate:
			http.Error(w, "Invalid creation date", http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(breed)
}

// GetCatBreed обрабатывает получение породы кошек по ID
func (h *CatBreedHandler) GetCatBreed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем ID из URL параметров
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	breed, err := h.service.GetCatBreedByID(id)
	if err != nil {
		http.Error(w, "Cat breed not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(breed)
}

// GetAllCatBreeds обрабатывает получение всех пород кошек
func (h *CatBreedHandler) GetAllCatBreeds(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	breeds, err := h.service.GetAllCatBreeds()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(breeds)
}

// GetUserCatBreeds обрабатывает получение пород кошек текущего пользователя
func (h *CatBreedHandler) GetUserCatBreeds(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем пользователя из контекста
	currentUser := middleware.GetUserFromContext(r.Context())
	if currentUser == nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	breeds, err := h.service.GetCatBreedsByUserID(currentUser.UserID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(breeds)
}

// UpdateCatBreed обрабатывает обновление породы кошек
func (h *CatBreedHandler) UpdateCatBreed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем пользователя из контекста
	currentUser := middleware.GetUserFromContext(r.Context())
	if currentUser == nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Извлекаем ID из URL параметров
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	var req models.CatBreedUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateCatBreed(id, &req, currentUser.UserID, currentUser.IsAdmin)
	if err != nil {
		switch err {
		case services.ErrCatBreedNotFound:
			http.Error(w, "Cat breed not found", http.StatusNotFound)
		case services.ErrCatBreedNameExists:
			http.Error(w, "Cat breed name already exists", http.StatusConflict)
		case services.ErrAccessDenied:
			http.Error(w, "Access denied", http.StatusForbidden)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Cat breed updated successfully"))
}

// DeleteCatBreed обрабатывает удаление породы кошек
func (h *CatBreedHandler) DeleteCatBreed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем пользователя из контекста
	currentUser := middleware.GetUserFromContext(r.Context())
	if currentUser == nil {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Извлекаем ID из URL параметров
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteCatBreed(id, currentUser.UserID, currentUser.IsAdmin)
	if err != nil {
		switch err {
		case services.ErrCatBreedNotFound:
			http.Error(w, "Cat breed not found", http.StatusNotFound)
		case services.ErrAccessDenied:
			http.Error(w, "Access denied", http.StatusForbidden)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Cat breed deleted successfully"))
}