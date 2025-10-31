package repositories

import (
	"meawle/internal/models"
)

// CatBreedRepository определяет интерфейс для работы с породами кошек
type CatBreedRepository interface {
	Create(breed *models.CatBreed) error
	GetByID(id int) (*models.CatBreed, error)
	GetAll() ([]models.CatBreed, error)
	GetByUserID(userID int) ([]models.CatBreed, error)
	Update(id int, breed *models.CatBreedUpdateRequest) error
	Delete(id int) error
	ExistsByName(name string) (bool, error)
	IsOwner(breedID int, userID int) (bool, error)
}

type catBreedRepository struct {
	db Database
}

// NewCatBreedRepository создает новый экземпляр репозитория пород кошек
func NewCatBreedRepository(db Database) CatBreedRepository {
	return &catBreedRepository{db: db}
}

// Create создает новую породу кошек
func (r *catBreedRepository) Create(breed *models.CatBreed) error {
	query := `INSERT INTO cat_breeds (name, description, user_id) VALUES (?, ?, ?)`

	result, err := r.db.Execute(query, breed.Name, breed.Description, breed.UserID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	breed.ID = int(id)
	return nil
}

// GetByID возвращает породу кошек по ID
func (r *catBreedRepository) GetByID(id int) (*models.CatBreed, error) {
	query := `SELECT id, name, description, user_id, created_at FROM cat_breeds WHERE id = ?`

	row := r.db.QueryRow(query, id)

	var breed models.CatBreed
	err := row.Scan(&breed.ID, &breed.Name, &breed.Description, &breed.UserID, &breed.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &breed, nil
}

// GetAll возвращает все породы кошек
func (r *catBreedRepository) GetAll() ([]models.CatBreed, error) {
	query := `SELECT id, name, description, user_id, created_at FROM cat_breeds ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var breeds []models.CatBreed
	for rows.Next() {
		var breed models.CatBreed
		err := rows.Scan(&breed.ID, &breed.Name, &breed.Description, &breed.UserID, &breed.CreatedAt)
		if err != nil {
			return nil, err
		}
		breeds = append(breeds, breed)
	}

	return breeds, nil
}

// GetByUserID возвращает породы кошек по ID пользователя
func (r *catBreedRepository) GetByUserID(userID int) ([]models.CatBreed, error) {
	query := `SELECT id, name, description, user_id, created_at FROM cat_breeds WHERE user_id = ? ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var breeds []models.CatBreed
	for rows.Next() {
		var breed models.CatBreed
		err := rows.Scan(&breed.ID, &breed.Name, &breed.Description, &breed.UserID, &breed.CreatedAt)
		if err != nil {
			return nil, err
		}
		breeds = append(breeds, breed)
	}

	return breeds, nil
}

// Update обновляет данные породы кошек
func (r *catBreedRepository) Update(id int, updateReq *models.CatBreedUpdateRequest) error {
	query := `UPDATE cat_breeds SET `
	params := []interface{}{}

	if updateReq.Name != nil {
		query += "name = ?, "
		params = append(params, *updateReq.Name)
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

// Delete удаляет породу кошек
func (r *catBreedRepository) Delete(id int) error {
	query := `DELETE FROM cat_breeds WHERE id = ?`
	_, err := r.db.Execute(query, id)
	return err
}

// ExistsByName проверяет существование породы по названию
func (r *catBreedRepository) ExistsByName(name string) (bool, error) {
	query := `SELECT COUNT(*) FROM cat_breeds WHERE name = ?`

	var count int
	err := r.db.QueryRow(query, name).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// IsOwner проверяет, является ли пользователь владельцем породы
func (r *catBreedRepository) IsOwner(breedID int, userID int) (bool, error) {
	query := `SELECT COUNT(*) FROM cat_breeds WHERE id = ? AND user_id = ?`

	var count int
	err := r.db.QueryRow(query, breedID, userID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}