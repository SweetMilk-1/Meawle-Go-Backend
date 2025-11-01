package repositories

import (
	"meawle/internal/models"
)

// CatRepository определяет интерфейс для работы с котами
type CatRepository interface {
	Create(cat *models.Cat) error
	GetByID(id int) (*models.Cat, error)
	GetAll() ([]models.Cat, error)
	GetByUserID(userID int) ([]models.Cat, error)
	Update(id int, cat *models.CatUpdateRequest) error
	Delete(id int) error
	IsOwner(catID int, userID int) (bool, error)
}

type catRepository struct {
	db Database
}

// NewCatRepository создает новый экземпляр репозитория котов
func NewCatRepository(db Database) CatRepository {
	return &catRepository{db: db}
}

// Create создает нового кота
func (r *catRepository) Create(cat *models.Cat) error {
	query := `INSERT INTO cats (name, age, description, user_id) VALUES (?, ?, ?, ?)`

	result, err := r.db.Execute(query, cat.Name, cat.Age, cat.Description, cat.UserID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	cat.ID = int(id)
	return nil
}

// GetByID возвращает кота по ID
func (r *catRepository) GetByID(id int) (*models.Cat, error) {
	query := `SELECT id, name, age, description, user_id, created_at FROM cats WHERE id = ?`

	row := r.db.QueryRow(query, id)

	var cat models.Cat
	err := row.Scan(&cat.ID, &cat.Name, &cat.Age, &cat.Description, &cat.UserID, &cat.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &cat, nil
}

// GetAll возвращает всех котов
func (r *catRepository) GetAll() ([]models.Cat, error) {
	query := `SELECT id, name, age, description, user_id, created_at FROM cats ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []models.Cat
	for rows.Next() {
		var cat models.Cat
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Age, &cat.Description, &cat.UserID, &cat.CreatedAt)
		if err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}

	return cats, nil
}

// GetByUserID возвращает котов по ID пользователя
func (r *catRepository) GetByUserID(userID int) ([]models.Cat, error) {
	query := `SELECT id, name, age, description, user_id, created_at FROM cats WHERE user_id = ? ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []models.Cat
	for rows.Next() {
		var cat models.Cat
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Age, &cat.Description, &cat.UserID, &cat.CreatedAt)
		if err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}

	return cats, nil
}

// Update обновляет данные кота
func (r *catRepository) Update(id int, updateReq *models.CatUpdateRequest) error {
	query := `UPDATE cats SET `
	params := []interface{}{}

	if updateReq.Name != nil {
		query += "name = ?, "
		params = append(params, *updateReq.Name)
	}

	if updateReq.Age != nil {
		query += "age = ?, "
		params = append(params, *updateReq.Age)
	}

	if updateReq.Description != nil {
		query += "description = ?, "
		params = append(params, *updateReq.Description)
	}

	// Убираем последнюю запятую и пробел
	query = query[:len(query)-2]
	query += " WHERE id = ?"
	params = append(params, id)

	_, err := r.db.Execute(query, params...)
	return err
}

// Delete удаляет кота
func (r *catRepository) Delete(id int) error {
	query := `DELETE FROM cats WHERE id = ?`
	_, err := r.db.Execute(query, id)
	return err
}

// IsOwner проверяет, является ли пользователь владельцем кота
func (r *catRepository) IsOwner(catID int, userID int) (bool, error) {
	query := `SELECT COUNT(*) FROM cats WHERE id = ? AND user_id = ?`

	var count int
	err := r.db.QueryRow(query, catID, userID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}