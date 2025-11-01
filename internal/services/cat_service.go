package services

import (
	"errors"
	"time"

	"meawle/internal/models"
	"meawle/internal/repositories"
)

var (
	ErrCatNotFound     = errors.New("cat not found")
	ErrInvalidCatData  = errors.New("invalid cat data")
	ErrInvalidCatAge   = errors.New("cat age must be between 0 and 30 years")
)

// CatService представляет сервис для работы с котами
type CatService struct {
	repo repositories.CatRepository
}

// NewCatService создает новый экземпляр сервиса котов
func NewCatService(repo repositories.CatRepository) *CatService {
	return &CatService{
		repo: repo,
	}
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

	// Создаем кота
	cat := &models.Cat{
		Name:        req.Name,
		Age:         req.Age,
		Description: req.Description,
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

// GetUserCats возвращает котов текущего пользователя
func (s *CatService) GetUserCats(userID int) ([]models.CatResponse, error) {
	cats, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []models.CatResponse
	for _, cat := range cats {
		responses = append(responses, cat.ToResponse())
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
	if !isAdmin && cat.UserID != userID {
		return ErrAccessDenied
	}

	// Проверяем возраст кота
	if req.Age != nil && (*req.Age < 0 || *req.Age > 30) {
		return ErrInvalidCatAge
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
	if !isAdmin && cat.UserID != userID {
		return ErrAccessDenied
	}

	return s.repo.Delete(id)
}