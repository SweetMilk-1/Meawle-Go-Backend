package services

import (
	"errors"
	"time"

	"meawle/internal/models"
	"meawle/internal/repositories"
)

var (
	ErrCatBreedNotFound    = errors.New("cat breed not found")
	ErrCatBreedNameExists  = errors.New("cat breed name already exists")
	ErrInvalidCatBreedData = errors.New("invalid cat breed data")
	ErrInvalidCreationDate = errors.New("creation date cannot be before 2000")
)

// CatBreedService представляет сервис для работы с породами кошек
type CatBreedService struct {
	repo repositories.CatBreedRepository
}

// NewCatBreedService создает новый экземпляр сервиса пород кошек
func NewCatBreedService(repo repositories.CatBreedRepository) *CatBreedService {
	return &CatBreedService{
		repo: repo,
	}
}

// Create создает новую породу кошек
func (s *CatBreedService) Create(req *models.CatBreedCreateRequest, userID int) (*models.CatBreedResponse, error) {
	// Проверяем валидацию названия
	if req.Name == "" {
		return nil, ErrInvalidCatBreedData
	}

	// Проверяем существование названия
	exists, err := s.repo.ExistsByName(req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrCatBreedNameExists
	}

	// Создаем породу
	breed := &models.CatBreed{
		Name:        req.Name,
		Description: req.Description,
		UserID:      userID,
		CreatedAt:   time.Now(),
	}

	// Проверяем дату создания (не должна быть раньше 2000 года)
	if breed.CreatedAt.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) {
		return nil, ErrInvalidCreationDate
	}

	err = s.repo.Create(breed)
	if err != nil {
		return nil, err
	}

	response := breed.ToResponse()
	return &response, nil
}

// GetCatBreedByID возвращает породу кошек по ID
func (s *CatBreedService) GetCatBreedByID(id int) (*models.CatBreedResponse, error) {
	breed, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrCatBreedNotFound
	}

	response := breed.ToResponse()
	return &response, nil
}

// GetAllCatBreeds возвращает все породы кошек
func (s *CatBreedService) GetAllCatBreeds() ([]models.CatBreedResponse, error) {
	breeds, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []models.CatBreedResponse
	for _, breed := range breeds {
		responses = append(responses, breed.ToResponse())
	}

	return responses, nil
}

// UpdateCatBreed обновляет данные породы кошек
func (s *CatBreedService) UpdateCatBreed(id int, req *models.CatBreedUpdateRequest, userID int, isAdmin bool) error {
	// Проверяем существование породы
	breed, err := s.repo.GetByID(id)
	if err != nil {
		return ErrCatBreedNotFound
	}

	// Проверяем права доступа: пользователь может обновлять только свои породы, админ - любые
	if !isAdmin && breed.UserID != userID {
		return ErrAccessDenied
	}

	// Если обновляется название, проверяем его уникальность
	if req.Name != nil {
		exists, err := s.repo.ExistsByName(*req.Name)
		if err != nil {
			return err
		}
		if exists && *req.Name != breed.Name {
			return ErrCatBreedNameExists
		}
	}

	return s.repo.Update(id, req)
}

// DeleteCatBreed удаляет породу кошек
func (s *CatBreedService) DeleteCatBreed(id int, userID int, isAdmin bool) error {
	// Проверяем существование породы
	breed, err := s.repo.GetByID(id)
	if err != nil {
		return ErrCatBreedNotFound
	}

	// Проверяем права доступа: пользователь может удалять только свои породы, админ - любые
	if !isAdmin && breed.UserID != userID {
		return ErrAccessDenied
	}

	return s.repo.Delete(id)
}
