package models

import (
	"time"
)

// CatBreed представляет модель породы кошек
type CatBreed struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// CatBreedCreateRequest представляет данные для создания породы кошек
type CatBreedCreateRequest struct {
	Name        string `json:"name" validate:"required,min=1"`
	Description string `json:"description" validate:"required,min=1"`
}

// CatBreedUpdateRequest представляет данные для обновления породы кошек
type CatBreedUpdateRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=1"`
}

// CatBreedResponse представляет ответ с данными породы кошек
type CatBreedResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// ToResponse преобразует CatBreed в CatBreedResponse
func (c *CatBreed) ToResponse() CatBreedResponse {
	return CatBreedResponse{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		UserID:      c.UserID,
		CreatedAt:   c.CreatedAt,
	}
}