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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.Error("Method not allowed"))
		return
	}

	// Получаем пользователя из контекста
	currentUser := middleware.GetUserFromContext(r.Context())
	if currentUser == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.Error("Authentication required"))
		return
	}

	var req models.CatBreedCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Error("Invalid request body"))
		return
	}

	breed, err := h.service.Create(&req, currentUser.UserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		switch err {
		case services.ErrCatBreedNameExists:
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(models.Error("Cat breed name already exists"))
		case services.ErrInvalidCatBreedData:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.Error("Invalid cat breed data"))
		case services.ErrInvalidCreationDate:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.Error("Invalid creation date"))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.Error("Internal server error"))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.Success(breed))
}

// GetCatBreed обрабатывает получение породы кошек по ID
func (h *CatBreedHandler) GetCatBreed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.Error("Method not allowed"))
		return
	}

	// Извлекаем ID из URL параметров
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Error("ID parameter is required"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Error("Invalid ID parameter"))
		return
	}

	breed, err := h.service.GetCatBreedByID(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.Error("Cat breed not found"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Success(breed))
}

// GetAllCatBreeds обрабатывает получение всех пород кошек
func (h *CatBreedHandler) GetAllCatBreeds(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.Error("Method not allowed"))
		return
	}

	breeds, err := h.service.GetAllCatBreeds()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error("Internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Success(breeds))
}

// GetUserCatBreeds обрабатывает получение пород кошек текущего пользователя
func (h *CatBreedHandler) GetUserCatBreeds(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.Error("Method not allowed"))
		return
	}

	// Получаем пользователя из контекста
	currentUser := middleware.GetUserFromContext(r.Context())
	if currentUser == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.Error("Authentication required"))
		return
	}

	breeds, err := h.service.GetCatBreedsByUserID(currentUser.UserID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.Error("Internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Success(breeds))
}

// UpdateCatBreed обрабатывает обновление породы кошек
func (h *CatBreedHandler) UpdateCatBreed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.Error("Method not allowed"))
		return
	}

	// Получаем пользователя из контекста
	currentUser := middleware.GetUserFromContext(r.Context())
	if currentUser == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.Error("Authentication required"))
		return
	}

	// Извлекаем ID из URL параметров
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Error("ID parameter is required"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Error("Invalid ID parameter"))
		return
	}

	var req models.CatBreedUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Error("Invalid request body"))
		return
	}

	err = h.service.UpdateCatBreed(id, &req, currentUser.UserID, currentUser.IsAdmin)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		switch err {
		case services.ErrCatBreedNotFound:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(models.Error("Cat breed not found"))
		case services.ErrCatBreedNameExists:
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(models.Error("Cat breed name already exists"))
		case services.ErrAccessDenied:
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(models.Error("Access denied"))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.Error("Internal server error"))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Success("Cat breed updated successfully"))
}

// DeleteCatBreed обрабатывает удаление породы кошек
func (h *CatBreedHandler) DeleteCatBreed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.Error("Method not allowed"))
		return
	}

	// Получаем пользователя из контекста
	currentUser := middleware.GetUserFromContext(r.Context())
	if currentUser == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.Error("Authentication required"))
		return
	}

	// Извлекаем ID из URL параметров
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Error("ID parameter is required"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.Error("Invalid ID parameter"))
		return
	}

	err = h.service.DeleteCatBreed(id, currentUser.UserID, currentUser.IsAdmin)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		switch err {
		case services.ErrCatBreedNotFound:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(models.Error("Cat breed not found"))
		case services.ErrAccessDenied:
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(models.Error("Access denied"))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.Error("Internal server error"))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Success("Cat breed deleted successfully"))
}
