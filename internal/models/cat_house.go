package models

import (
	"time"
)

// CatHouse представляет модель домика для кошек
type CatHouse struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Capacity         int       `json:"capacity"`
	CurrentOccupancy int       `json:"current_occupancy"`
	UserID           int       `json:"user_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// CatHouseCreateRequest представляет данные для создания домика
type CatHouseCreateRequest struct {
	Name     string `json:"name" validate:"required,min=1"`
	Capacity int    `json:"capacity" validate:"required,min=1"`
}

// CatHouseUpdateRequest представляет данные для обновления домика
type CatHouseUpdateRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=1"`
	Capacity *int    `json:"capacity,omitempty" validate:"omitempty,min=1"`
}

// CatHouseResponse представляет ответ с данными домика
type CatHouseResponse struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Capacity         int       `json:"capacity"`
	CurrentOccupancy int       `json:"current_occupancy"`
	AvailableSpaces  int       `json:"available_spaces"`
	UserID           int       `json:"user_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	CanEdit          bool      `json:"can_edit"`
}

// CatHouseOccupancy представляет связь кота с домиком
type CatHouseOccupancy struct {
	ID         int       `json:"id"`
	CatID      int       `json:"cat_id"`
	CatHouseID int       `json:"cat_house_id"`
	OccupiedAt time.Time `json:"occupied_at"`
}

// CatHouseOccupancyRequest представляет запрос на размещение кота в домике
type CatHouseOccupancyRequest struct {
	CatHouseID int `json:"cat_house_id" validate:"required,min=1"`
}

// CatHouseOccupancyResponse представляет ответ с информацией о размещении
type CatHouseOccupancyResponse struct {
	ID           int       `json:"id"`
	CatID        int       `json:"cat_id"`
	CatHouseID   int       `json:"cat_house_id"`
	CatName      string    `json:"cat_name"`
	CatHouseName string    `json:"cat_house_name"`
	OccupiedAt   time.Time `json:"occupied_at"`
}

// ToResponse преобразует CatHouse в CatHouseResponse
func (ch *CatHouse) ToResponse() CatHouseResponse {
	return CatHouseResponse{
		ID:               ch.ID,
		Name:             ch.Name,
		Capacity:         ch.Capacity,
		CurrentOccupancy: ch.CurrentOccupancy,
		AvailableSpaces:  ch.Capacity - ch.CurrentOccupancy,
		UserID:           ch.UserID,
		CreatedAt:        ch.CreatedAt,
		UpdatedAt:        ch.UpdatedAt,
		CanEdit:          false, // По умолчанию false, будет установлено в сервисе
	}
}

// ToResponse преобразует CatHouseOccupancy в CatHouseOccupancyResponse
func (cho *CatHouseOccupancy) ToResponse(catName, catHouseName string) CatHouseOccupancyResponse {
	return CatHouseOccupancyResponse{
		ID:           cho.ID,
		CatID:        cho.CatID,
		CatHouseID:   cho.CatHouseID,
		CatName:      catName,
		CatHouseName: catHouseName,
		OccupiedAt:   cho.OccupiedAt,
	}
}
