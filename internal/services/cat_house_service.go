package services

import (
	"errors"
	"time"

	"meawle/internal/models"
	"meawle/internal/repositories"
)

var (
	ErrCatHouseNotFound     = errors.New("cat house not found")
	ErrInvalidCatHouseData  = errors.New("invalid cat house data")
	ErrCatHouseFull         = errors.New("cat house is full")
	ErrCatAlreadyInHouse    = errors.New("cat is already in a house")
	ErrCatNotInHouse        = errors.New("cat is not in a house")
	ErrCatNotFoundInHouse   = errors.New("cat not found in house")
	ErrInsufficientCapacity = errors.New("new capacity is less than current occupancy")
)

// CatHouseService представляет сервис для работы с домиками для кошек
type CatHouseService struct {
	repo    repositories.CatHouseRepository
	catRepo repositories.CatRepository
}

// NewCatHouseService создает новый экземпляр сервиса домиков для кошек
func NewCatHouseService(repo repositories.CatHouseRepository, catRepo repositories.CatRepository) *CatHouseService {
	return &CatHouseService{
		repo:    repo,
		catRepo: catRepo,
	}
}

// isEditableByUser проверяет, может ли пользователь редактировать домик
func (s *CatHouseService) isEditableByUser(houseUserID, currentUserID int, isAdmin bool) bool {
	// Если userID = 0 - неавторизованный пользователь, can_edit = false
	if currentUserID == 0 {
		return false
	}
	return isAdmin || houseUserID == currentUserID
}

// Create создает новый домик
func (s *CatHouseService) Create(req *models.CatHouseCreateRequest, userID int) (*models.CatHouseResponse, error) {
	// Проверяем валидацию названия и вместимости
	if req.Name == "" || req.Capacity <= 0 {
		return nil, ErrInvalidCatHouseData
	}

	// Создаем домик
	house := &models.CatHouse{
		Name:             req.Name,
		Capacity:         req.Capacity,
		CurrentOccupancy: 0,
		UserID:           userID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := s.repo.Create(house)
	if err != nil {
		return nil, err
	}

	response := house.ToResponse()
	response.CanEdit = true // Пользователь всегда может редактировать свой домик
	return &response, nil
}

// GetCatHouseByID возвращает домик по ID
func (s *CatHouseService) GetCatHouseByID(id int) (*models.CatHouseResponse, error) {
	house, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrCatHouseNotFound
	}

	response := house.ToResponse()
	return &response, nil
}

// GetCatHouseByIDWithUser возвращает домик по ID с учетом прав доступа пользователя
func (s *CatHouseService) GetCatHouseByIDWithUser(id int, currentUserID int, isAdmin bool) (*models.CatHouseResponse, error) {
	house, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrCatHouseNotFound
	}

	response := house.ToResponse()
	response.CanEdit = s.isEditableByUser(house.UserID, currentUserID, isAdmin)
	return &response, nil
}

// GetAllCatHouses возвращает все домики
func (s *CatHouseService) GetAllCatHouses() ([]models.CatHouseResponse, error) {
	houses, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []models.CatHouseResponse
	for _, house := range houses {
		responses = append(responses, house.ToResponse())
	}

	return responses, nil
}

// GetAllCatHousesWithUser возвращает все домики с учетом прав доступа пользователя
func (s *CatHouseService) GetAllCatHousesWithUser(currentUserID int, isAdmin bool) ([]models.CatHouseResponse, error) {
	houses, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []models.CatHouseResponse
	for _, house := range houses {
		response := house.ToResponse()
		response.CanEdit = s.isEditableByUser(house.UserID, currentUserID, isAdmin)
		responses = append(responses, response)
	}

	return responses, nil
}

// UpdateCatHouse обновляет данные домика
func (s *CatHouseService) UpdateCatHouse(id int, req *models.CatHouseUpdateRequest, userID int, isAdmin bool) error {
	// Проверяем существование домика
	house, err := s.repo.GetByID(id)
	if err != nil {
		return ErrCatHouseNotFound
	}

	// Проверяем права доступа: пользователь может обновлять только свои домики, админ - любые
	if !s.isEditableByUser(house.UserID, userID, isAdmin) {
		return ErrAccessDenied
	}

	// Проверяем, что новая вместимость не меньше текущего количества котов
	if req.Capacity != nil && *req.Capacity < house.CurrentOccupancy {
		return ErrInsufficientCapacity
	}

	return s.repo.Update(id, req)
}

// DeleteCatHouse удаляет домик
func (s *CatHouseService) DeleteCatHouse(id int, userID int, isAdmin bool) error {
	// Проверяем существование домика
	house, err := s.repo.GetByID(id)
	if err != nil {
		return ErrCatHouseNotFound
	}

	// Проверяем права доступа: пользователь может удалять только свои домики, админ - любые
	if !s.isEditableByUser(house.UserID, userID, isAdmin) {
		return ErrAccessDenied
	}

	return s.repo.Delete(id)
}

// AddCatToHouse добавляет кота в домик
func (s *CatHouseService) AddCatToHouse(catID, houseID int) (*models.CatHouseOccupancyResponse, error) {
	// Проверяем существование кота
	cat, err := s.catRepo.GetByID(catID)
	if err != nil {
		return nil, ErrCatNotFound
	}

	// Проверяем существование домика
	house, err := s.repo.GetByID(houseID)
	if err != nil {
		return nil, ErrCatHouseNotFound
	}

	// Проверяем, есть ли место в домике
	if house.CurrentOccupancy >= house.Capacity {
		return nil, ErrCatHouseFull
	}

	// Проверяем, не находится ли кот уже в каком-то домике
	existingOccupancy, err := s.repo.GetOccupancyByCatID(catID)
	if err == nil && existingOccupancy != nil {
		return nil, ErrCatAlreadyInHouse
	}

	// Добавляем кота в домик
	err = s.repo.AddCatToHouse(catID, houseID)
	if err != nil {
		return nil, err
	}

	// Получаем обновленную информацию о размещении
	occupancy, err := s.repo.GetOccupancyByCatID(catID)
	if err != nil {
		return nil, err
	}

	response := occupancy.ToResponse(cat.Name, house.Name)
	return &response, nil
}

// RemoveCatFromHouse удаляет кота из домика
func (s *CatHouseService) RemoveCatFromHouse(catID int) error {
	// Проверяем существование кота
	_, err := s.catRepo.GetByID(catID)
	if err != nil {
		return ErrCatNotFound
	}

	// Проверяем, находится ли кот в каком-то домике
	existingOccupancy, err := s.repo.GetOccupancyByCatID(catID)
	if err != nil || existingOccupancy == nil {
		return ErrCatNotInHouse
	}

	return s.repo.RemoveCatFromHouse(catID)
}

// GetCatHouseWithCats возвращает домик со списком котов в нем
func (s *CatHouseService) GetCatHouseWithCats(houseID int) (*models.CatHouseResponse, []models.CatResponse, error) {
	// Получаем домик с котами
	house, cats, err := s.repo.GetHouseWithCats(houseID)
	if err != nil {
		return nil, nil, ErrCatHouseNotFound
	}

	// Преобразуем домик в ответ
	houseResponse := house.ToResponse()

	// Преобразуем котов в ответы
	var catResponses []models.CatResponse
	for _, cat := range cats {
		catResponses = append(catResponses, cat.ToResponse())
	}

	return &houseResponse, catResponses, nil
}

// GetCatHouseWithCatsWithUser возвращает домик со списком котов в нем с учетом прав доступа пользователя
func (s *CatHouseService) GetCatHouseWithCatsWithUser(houseID int, currentUserID int, isAdmin bool) (*models.CatHouseResponse, []models.CatResponse, error) {
	// Получаем домик с котами
	house, cats, err := s.repo.GetHouseWithCats(houseID)
	if err != nil {
		return nil, nil, ErrCatHouseNotFound
	}

	// Преобразуем домик в ответ
	houseResponse := house.ToResponse()
	houseResponse.CanEdit = s.isEditableByUser(house.UserID, currentUserID, isAdmin)

	// Преобразуем котов в ответы с учетом прав доступа
	var catResponses []models.CatResponse
	for _, cat := range cats {
		response := cat.ToResponse()
		// Здесь нужно установить can_edit для котов, но у нас нет доступа к сервису котов
		// Временно оставляем false, можно будет добавить позже
		catResponses = append(catResponses, response)
	}

	return &houseResponse, catResponses, nil
}

// GetCatHouseOccupancy возвращает информацию о размещении конкретного кота
func (s *CatHouseService) GetCatHouseOccupancy(catID int) (*models.CatHouseOccupancyResponse, error) {
	// Проверяем существование кота
	cat, err := s.catRepo.GetByID(catID)
	if err != nil {
		return nil, ErrCatNotFound
	}

	// Получаем информацию о размещении
	occupancy, err := s.repo.GetOccupancyByCatID(catID)
	if err != nil || occupancy == nil {
		return nil, ErrCatNotInHouse
	}

	// Получаем информацию о домике
	house, err := s.repo.GetByID(occupancy.CatHouseID)
	if err != nil {
		return nil, ErrCatHouseNotFound
	}

	response := occupancy.ToResponse(cat.Name, house.Name)
	return &response, nil
}

// GetCatsInHouse возвращает список котов в домике
func (s *CatHouseService) GetCatsInHouse(houseID int) ([]models.CatHouseOccupancyResponse, error) {
	// Проверяем существование домика
	house, err := s.repo.GetByID(houseID)
	if err != nil {
		return nil, ErrCatHouseNotFound
	}

	// Получаем информацию о размещении
	occupancies, err := s.repo.GetOccupancyByHouseID(houseID)
	if err != nil {
		return nil, err
	}

	var responses []models.CatHouseOccupancyResponse
	for _, occupancy := range occupancies {
		// Получаем информацию о коте
		cat, err := s.catRepo.GetByID(occupancy.CatID)
		if err != nil {
			continue // Пропускаем, если кот не найден
		}

		response := occupancy.ToResponse(cat.Name, house.Name)
		responses = append(responses, response)
	}

	return responses, nil
}
