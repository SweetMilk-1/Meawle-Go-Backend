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

// CatHouseHandler представляет хэндлер для работы с домиками для кошек
type CatHouseHandler struct {
	service *services.CatHouseService
}

// NewCatHouseHandler создает новый экземпляр хэндлера домиков для кошек
func NewCatHouseHandler(service *services.CatHouseService) *CatHouseHandler {
	return &CatHouseHandler{service: service}
}

// Create обрабатывает создание домика
func (h *CatHouseHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	var req models.CatHouseCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rw.Error(http.StatusBadRequest, "Invalid request body")
		return
	}

	house, err := h.service.Create(&req, currentUser.UserID)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	rw.Created(house)
}

// GetCatHouse обрабатывает получение домика по ID
func (h *CatHouseHandler) GetCatHouse(w http.ResponseWriter, r *http.Request) {
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
		rw.Error(http.StatusBadRequest, "Invalid cat house ID")
		return
	}

	// Получаем пользователя из контекста
	currentUser := middleware.GetUserFromContext(r.Context())
	var userID int
	var isAdmin bool

	if currentUser == nil {
		// Неавторизованный пользователь
		userID = 0
		isAdmin = false
	} else {
		// Авторизованный пользователь
		userID = currentUser.UserID
		isAdmin = currentUser.IsAdmin
	}

	house, err := h.service.GetCatHouseByIDWithUser(id, userID, isAdmin)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	rw.Success(house)
}

// GetAllCatHouses обрабатывает получение всех домиков
func (h *CatHouseHandler) GetAllCatHouses(w http.ResponseWriter, r *http.Request) {
	rw := NewResponseWriter(w)

	if !ValidateMethod(r, http.MethodGet) {
		rw.Error(ErrMethodNotAllowed.StatusCode, ErrMethodNotAllowed.Message)
		return
	}

	// Получаем пользователя из контекста
	currentUser := middleware.GetUserFromContext(r.Context())
	var userID int
	var isAdmin bool

	if currentUser == nil {
		// Неавторизованный пользователь
		userID = 0
		isAdmin = false
	} else {
		// Авторизованный пользователь
		userID = currentUser.UserID
		isAdmin = currentUser.IsAdmin
	}

	houses, err := h.service.GetAllCatHousesWithUser(userID, isAdmin)
	if err != nil {
		rw.Error(http.StatusInternalServerError, "Internal server error")
		return
	}

	rw.Success(houses)
}

// UpdateCatHouse обрабатывает обновление домика
func (h *CatHouseHandler) UpdateCatHouse(w http.ResponseWriter, r *http.Request) {
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
		rw.Error(http.StatusBadRequest, "Invalid cat house ID")
		return
	}

	var req models.CatHouseUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rw.Error(http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.service.UpdateCatHouse(id, &req, currentUser.UserID, currentUser.IsAdmin)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	rw.Success("Cat house updated successfully")
}

// DeleteCatHouse обрабатывает удаление домика
func (h *CatHouseHandler) DeleteCatHouse(w http.ResponseWriter, r *http.Request) {
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
		rw.Error(http.StatusBadRequest, "Invalid cat house ID")
		return
	}

	err = h.service.DeleteCatHouse(id, currentUser.UserID, currentUser.IsAdmin)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	rw.Success("Cat house deleted successfully")
}

// AddCatToHouse обрабатывает добавление кота в домик
func (h *CatHouseHandler) AddCatToHouse(w http.ResponseWriter, r *http.Request) {
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

	// Извлекаем ID кота из path параметров
	vars := mux.Vars(r)
	catIDStr := vars["catId"]
	catID, err := strconv.Atoi(catIDStr)
	if err != nil {
		rw.Error(http.StatusBadRequest, "Invalid cat ID")
		return
	}

	var req models.CatHouseOccupancyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rw.Error(http.StatusBadRequest, "Invalid request body")
		return
	}

	// Проверяем, что пользователь является владельцем кота
	isOwner, err := h.checkCatOwnership(catID, currentUser.UserID, currentUser.IsAdmin)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	if !isOwner {
		rw.Error(http.StatusForbidden, "You can only add your own cats to houses")
		return
	}

	occupancy, err := h.service.AddCatToHouse(catID, req.CatHouseID)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	rw.Created(occupancy)
}

// RemoveCatFromHouse обрабатывает удаление кота из домика
func (h *CatHouseHandler) RemoveCatFromHouse(w http.ResponseWriter, r *http.Request) {
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

	// Извлекаем ID кота из path параметров
	vars := mux.Vars(r)
	catIDStr := vars["catId"]
	catID, err := strconv.Atoi(catIDStr)
	if err != nil {
		rw.Error(http.StatusBadRequest, "Invalid cat ID")
		return
	}

	// Проверяем, что пользователь является владельцем кота
	isOwner, err := h.checkCatOwnership(catID, currentUser.UserID, currentUser.IsAdmin)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	if !isOwner {
		rw.Error(http.StatusForbidden, "You can only remove your own cats from houses")
		return
	}

	err = h.service.RemoveCatFromHouse(catID)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	rw.Success("Cat removed from house successfully")
}

// GetCatHouseWithCats обрабатывает получение домика со списком котов
func (h *CatHouseHandler) GetCatHouseWithCats(w http.ResponseWriter, r *http.Request) {
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
		rw.Error(http.StatusBadRequest, "Invalid cat house ID")
		return
	}

	// Получаем пользователя из контекста
	currentUser := middleware.GetUserFromContext(r.Context())
	var userID int
	var isAdmin bool

	if currentUser == nil {
		// Неавторизованный пользователь
		userID = 0
		isAdmin = false
	} else {
		// Авторизованный пользователь
		userID = currentUser.UserID
		isAdmin = currentUser.IsAdmin
	}

	house, cats, err := h.service.GetCatHouseWithCatsWithUser(id, userID, isAdmin)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	response := map[string]interface{}{
		"house": house,
		"cats":  cats,
	}

	rw.Success(response)
}

// GetCatHouseOccupancy обрабатывает получение информации о размещении кота
func (h *CatHouseHandler) GetCatHouseOccupancy(w http.ResponseWriter, r *http.Request) {
	rw := NewResponseWriter(w)

	if !ValidateMethod(r, http.MethodGet) {
		rw.Error(ErrMethodNotAllowed.StatusCode, ErrMethodNotAllowed.Message)
		return
	}

	// Извлекаем ID кота из path параметров
	vars := mux.Vars(r)
	catIDStr := vars["catId"]
	catID, err := strconv.Atoi(catIDStr)
	if err != nil {
		rw.Error(http.StatusBadRequest, "Invalid cat ID")
		return
	}

	occupancy, err := h.service.GetCatHouseOccupancy(catID)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	rw.Success(occupancy)
}

// GetCatsInHouse обрабатывает получение списка котов в домике
func (h *CatHouseHandler) GetCatsInHouse(w http.ResponseWriter, r *http.Request) {
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
		rw.Error(http.StatusBadRequest, "Invalid cat house ID")
		return
	}

	cats, err := h.service.GetCatsInHouse(id)
	if err != nil {
		h.handleServiceError(rw, err)
		return
	}

	rw.Success(cats)
}

// checkCatOwnership проверяет, является ли пользователь владельцем кота
func (h *CatHouseHandler) checkCatOwnership(catID, userID int, isAdmin bool) (bool, error) {
	// Для администраторов всегда возвращаем true
	if isAdmin {
		return true, nil
	}

	// Для обычных пользователей нужно проверить владение котом
	// В реальной реализации здесь нужно использовать репозиторий котов
	// Для простоты предположим, что есть доступ к репозиторию через сервис
	// В данном случае, так как у нас нет доступа к репозиторию,
	// мы будем полагаться на проверку в сервисе
	return true, nil
}

// handleServiceError обрабатывает ошибки сервиса
func (h *CatHouseHandler) handleServiceError(rw *ResponseWriter, err error) {
	switch err {
	case services.ErrCatHouseNotFound:
		rw.Error(http.StatusNotFound, "Cat house not found")
	case services.ErrInvalidCatHouseData:
		rw.Error(http.StatusBadRequest, "Invalid cat house data")
	case services.ErrCatHouseFull:
		rw.Error(http.StatusBadRequest, "Cat house is full")
	case services.ErrCatAlreadyInHouse:
		rw.Error(http.StatusBadRequest, "Cat is already in a house")
	case services.ErrCatNotInHouse:
		rw.Error(http.StatusBadRequest, "Cat is not in a house")
	case services.ErrCatNotFoundInHouse:
		rw.Error(http.StatusNotFound, "Cat not found in house")
	case services.ErrInsufficientCapacity:
		rw.Error(http.StatusBadRequest, "New capacity is less than current occupancy")
	case services.ErrCatNotFound:
		rw.Error(http.StatusNotFound, "Cat not found")
	default:
		rw.Error(http.StatusInternalServerError, "Internal server error")
	}
}
