package models

import (
	"time"
)

// Cat представляет модель кота
type Cat struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Age         *int      `json:"age,omitempty"`
	Description *string   `json:"description,omitempty"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// CatCreateRequest представляет данные для создания кота
type CatCreateRequest struct {
	Name        string  `json:"name" validate:"required,min=1"`
	Age         *int    `json:"age,omitempty" validate:"omitempty,min=0,max=30"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=1"`
}

// CatUpdateRequest представляет данные для обновления кота
type CatUpdateRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1"`
	Age         *int    `json:"age,omitempty" validate:"omitempty,min=0,max=30"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=1"`
}

// CatResponse представляет ответ с данными кота
type CatResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Age         *int      `json:"age,omitempty"`
	Description *string   `json:"description,omitempty"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

// ToResponse преобразует Cat в CatResponse
func (c *Cat) ToResponse() CatResponse {
	return CatResponse{
		ID:          c.ID,
		Name:        c.Name,
		Age:         c.Age,
		Description: c.Description,
		UserID:      c.UserID,
		CreatedAt:   c.CreatedAt,
	}
}