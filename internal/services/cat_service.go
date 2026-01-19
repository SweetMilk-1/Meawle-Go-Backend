package services

import (
	"errors"
	"time"

	"meawle/internal/models"
	"meawle/internal/repositories"
)

var (
	ErrCatNotFound    = errors.New("cat not found")
	ErrInvalidCatData = errors.New("invalid cat data")
	ErrInvalidCatAge  = errors.New("cat age must be between 0 and 30 years")
	ErrAccessDenied   = errors.New("access denied")
)

// CatService представляет сервис для работы с котами
type CatService struct {
	repo      repositories.CatRepository
	breedRepo repositories.CatBreedRepository
}

// NewCatService создает новый экземпляр сервиса котов
func NewCatService(repo repositories.CatRepository, breedRepo repositories.CatBreedRepository) *CatService {
	return &CatService{
		repo:      repo,
		breedRepo: breedRepo,
	}
}

// isEditableByUser проверяет, может ли пользователь редактировать кота
func (s *CatService) isEditableByUser(catUserID, currentUserID int, isAdmin bool) bool {
	// Если userID = 0 - неавторизованный пользователь, can_edit = false
	if currentUserID == 0 {
		return false
	}
	return isAdmin || catUserID == currentUserID
}

// checkCatBreedExists проверяет существование породы кошки
func (s *CatService) checkCatBreedExists(catBreedID *int) error {
	if catBreedID == nil {
		return nil // Поле опциональное, nil допустимо
	}

	_, err := s.breedRepo.GetByID(*catBreedID)
	if err != nil {
		return ErrCatBreedNotFound // Используем ошибку из cat_breed_service.go
	}

	return nil
}

// Create создает нового кота
func (s *CatService) Create(req *models.CatCreateRequest, userID int) (*models.CatResponse, error) {
	// Проверяем валидацию названия
	if req.Name == "" {
		return nil, ErrInvalidCatData
	}

	// Проверяем возраст кота
	if req.Age != nil && (*req.Age < 0 || *req.Age > 30) {
		return nil, ErrInvalidCatAge
	}

	// Проверяем существование породы, если указана
	if err := s.checkCatBreedExists(req.CatBreedID); err != nil {
		return nil, err
	}

	// Создаем кота
	cat := &models.Cat{
		Name:        req.Name,
		Age:         req.Age,
		Description: req.Description,
		CatBreedID:  req.CatBreedID,
		UserID:      userID,
		CreatedAt:   time.Now(),
	}

	err := s.repo.Create(cat)
	if err != nil {
		return nil, err
	}

	response := cat.ToResponse()
	return &response, nil
}

// GetCatByID возвращает кота по ID
func (s *CatService) GetCatByID(id int) (*models.CatResponse, error) {
	cat, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrCatNotFound
	}

	response := cat.ToResponse()
	return &response, nil
}

// GetCatByIDWithUser возвращает кота по ID с учетом прав доступа пользователя
func (s *CatService) GetCatByIDWithUser(id int, currentUserID int, isAdmin bool) (*models.CatResponse, error) {
	cat, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrCatNotFound
	}

	response := cat.ToResponse()
	response.CanEdit = s.isEditableByUser(cat.UserID, currentUserID, isAdmin)
	return &response, nil
}

// GetAllCats возвращает всех котов
func (s *CatService) GetAllCats() ([]models.CatResponse, error) {
	cats, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []models.CatResponse
	for _, cat := range cats {
		responses = append(responses, cat.ToResponse())
	}

	return responses, nil
}

// GetAllCatsWithUser возвращает всех котов с учетом прав доступа пользователя
func (s *CatService) GetAllCatsWithUser(currentUserID int, isAdmin bool) ([]models.CatResponse, error) {
	cats, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []models.CatResponse
	for _, cat := range cats {
		response := cat.ToResponse()
		response.CanEdit = s.isEditableByUser(cat.UserID, currentUserID, isAdmin)
		responses = append(responses, response)
	}

	return responses, nil
}

// GetUserCats возвращает котов текущего пользователя
func (s *CatService) GetUserCats(userID int) ([]models.CatResponse, error) {
	cats, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []models.CatResponse
	for _, cat := range cats {
		response := cat.ToResponse()
		response.CanEdit = true // Пользователь всегда может редактировать своих котов
		responses = append(responses, response)
	}

	return responses, nil
}

// UpdateCat обновляет данные кота
func (s *CatService) UpdateCat(id int, req *models.CatUpdateRequest, userID int, isAdmin bool) error {
	// Проверяем существование кота
	cat, err := s.repo.GetByID(id)
	if err != nil {
		return ErrCatNotFound
	}

	// Проверяем права доступа: пользователь может обновлять только своих котов, админ - любых
	if !s.isEditableByUser(cat.UserID, userID, isAdmin) {
		return ErrAccessDenied
	}

	// Проверяем возраст кота
	if req.Age != nil && (*req.Age < 0 || *req.Age > 30) {
		return ErrInvalidCatAge
	}

	// Проверяем существование породы, если указана
	if err := s.checkCatBreedExists(req.CatBreedID); err != nil {
		return err
	}

	return s.repo.Update(id, req)
}

// DeleteCat удаляет кота
func (s *CatService) DeleteCat(id int, userID int, isAdmin bool) error {
	// Проверяем существование кота
	cat, err := s.repo.GetByID(id)
	if err != nil {
		return ErrCatNotFound
	}

	// Проверяем права доступа: пользователь может удалять только своих котов, админ - любых
	if !s.isEditableByUser(cat.UserID, userID, isAdmin) {
		return ErrAccessDenied
	}

	return s.repo.Delete(id)
}
