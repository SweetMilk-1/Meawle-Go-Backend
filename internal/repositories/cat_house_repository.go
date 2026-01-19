package repositories

import (
	"meawle/internal/models"
)

// CatHouseRepository определяет интерфейс для работы с домиками для кошек
type CatHouseRepository interface {
	Create(house *models.CatHouse) error
	GetByID(id int) (*models.CatHouse, error)
	GetAll() ([]models.CatHouse, error)
	Update(id int, house *models.CatHouseUpdateRequest) error
	Delete(id int) error
	GetOccupancyByCatID(catID int) (*models.CatHouseOccupancy, error)
	GetOccupancyByHouseID(houseID int) ([]models.CatHouseOccupancy, error)
	AddCatToHouse(catID, houseID int) error
	RemoveCatFromHouse(catID int) error
	UpdateHouseOccupancy(houseID int) error
	GetHouseWithCats(houseID int) (*models.CatHouse, []models.Cat, error)
}

type catHouseRepository struct {
	db Database
}

// NewCatHouseRepository создает новый экземпляр репозитория домиков для кошек
func NewCatHouseRepository(db Database) CatHouseRepository {
	return &catHouseRepository{db: db}
}

// Create создает новый домик
func (r *catHouseRepository) Create(house *models.CatHouse) error {
	query := `INSERT INTO cat_houses (name, capacity, current_occupancy, user_id) VALUES (?, ?, ?, ?)`

	result, err := r.db.Execute(query, house.Name, house.Capacity, house.CurrentOccupancy, house.UserID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	house.ID = int(id)
	return nil
}

// GetByID возвращает домик по ID
func (r *catHouseRepository) GetByID(id int) (*models.CatHouse, error) {
	query := `SELECT id, name, capacity, current_occupancy, user_id, created_at, updated_at FROM cat_houses WHERE id = ?`

	row := r.db.QueryRow(query, id)

	var house models.CatHouse
	err := row.Scan(&house.ID, &house.Name, &house.Capacity, &house.CurrentOccupancy, &house.UserID, &house.CreatedAt, &house.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &house, nil
}

// GetAll возвращает все домики
func (r *catHouseRepository) GetAll() ([]models.CatHouse, error) {
	query := `SELECT id, name, capacity, current_occupancy, user_id, created_at, updated_at FROM cat_houses ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var houses []models.CatHouse
	for rows.Next() {
		var house models.CatHouse
		err := rows.Scan(&house.ID, &house.Name, &house.Capacity, &house.CurrentOccupancy, &house.UserID, &house.CreatedAt, &house.UpdatedAt)
		if err != nil {
			return nil, err
		}
		houses = append(houses, house)
	}

	return houses, nil
}

// Update обновляет данные домика
func (r *catHouseRepository) Update(id int, updateReq *models.CatHouseUpdateRequest) error {
	query := `UPDATE cat_houses SET `
	params := []interface{}{}

	if updateReq.Name != nil {
		query += "name = ?, "
		params = append(params, *updateReq.Name)
	}

	if updateReq.Capacity != nil {
		query += "capacity = ?, "
		params = append(params, *updateReq.Capacity)
	}

	// Добавляем обновление времени
	query += "updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	params = append(params, id)

	_, err := r.db.Execute(query, params...)
	return err
}

// Delete удаляет домик
func (r *catHouseRepository) Delete(id int) error {
	// Начинаем транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Удаляем все связи с котами
	_, err = tx.Exec("DELETE FROM cat_house_occupancy WHERE cat_house_id = ?", id)
	if err != nil {
		return err
	}

	// Удаляем домик
	_, err = tx.Exec("DELETE FROM cat_houses WHERE id = ?", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetOccupancyByCatID возвращает информацию о размещении кота
func (r *catHouseRepository) GetOccupancyByCatID(catID int) (*models.CatHouseOccupancy, error) {
	query := `SELECT id, cat_id, cat_house_id, occupied_at FROM cat_house_occupancy WHERE cat_id = ?`

	row := r.db.QueryRow(query, catID)

	var occupancy models.CatHouseOccupancy
	err := row.Scan(&occupancy.ID, &occupancy.CatID, &occupancy.CatHouseID, &occupancy.OccupiedAt)
	if err != nil {
		return nil, err
	}

	return &occupancy, nil
}

// GetOccupancyByHouseID возвращает всех котов в домике
func (r *catHouseRepository) GetOccupancyByHouseID(houseID int) ([]models.CatHouseOccupancy, error) {
	query := `SELECT id, cat_id, cat_house_id, occupied_at FROM cat_house_occupancy WHERE cat_house_id = ? ORDER BY occupied_at DESC`

	rows, err := r.db.Query(query, houseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var occupancies []models.CatHouseOccupancy
	for rows.Next() {
		var occupancy models.CatHouseOccupancy
		err := rows.Scan(&occupancy.ID, &occupancy.CatID, &occupancy.CatHouseID, &occupancy.OccupiedAt)
		if err != nil {
			return nil, err
		}
		occupancies = append(occupancies, occupancy)
	}

	return occupancies, nil
}

// AddCatToHouse добавляет кота в домик
func (r *catHouseRepository) AddCatToHouse(catID, houseID int) error {
	// Начинаем транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Добавляем связь кота с домиком
	query := `INSERT INTO cat_house_occupancy (cat_id, cat_house_id) VALUES (?, ?)`
	_, err = tx.Exec(query, catID, houseID)
	if err != nil {
		return err
	}

	// Обновляем количество занятых мест в домике
	updateQuery := `UPDATE cat_houses SET current_occupancy = current_occupancy + 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err = tx.Exec(updateQuery, houseID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// RemoveCatFromHouse удаляет кота из домика
func (r *catHouseRepository) RemoveCatFromHouse(catID int) error {
	// Начинаем транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Получаем ID домика, из которого удаляем кота
	var houseID int
	err = tx.QueryRow("SELECT cat_house_id FROM cat_house_occupancy WHERE cat_id = ?", catID).Scan(&houseID)
	if err != nil {
		return err
	}

	// Удаляем связь кота с домиком
	_, err = tx.Exec("DELETE FROM cat_house_occupancy WHERE cat_id = ?", catID)
	if err != nil {
		return err
	}

	// Обновляем количество занятых мест в домике
	updateQuery := `UPDATE cat_houses SET current_occupancy = current_occupancy - 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err = tx.Exec(updateQuery, houseID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// UpdateHouseOccupancy обновляет количество занятых мест в домике
func (r *catHouseRepository) UpdateHouseOccupancy(houseID int) error {
	query := `
		UPDATE cat_houses 
		SET current_occupancy = (
			SELECT COUNT(*) FROM cat_house_occupancy WHERE cat_house_id = ?
		), 
		updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`

	_, err := r.db.Execute(query, houseID, houseID)
	return err
}

// GetHouseWithCats возвращает домик со списком котов в нем
func (r *catHouseRepository) GetHouseWithCats(houseID int) (*models.CatHouse, []models.Cat, error) {
	// Получаем домик
	house, err := r.GetByID(houseID)
	if err != nil {
		return nil, nil, err
	}

	// Получаем котов в домике
	query := `
		SELECT c.id, c.name, c.age, c.description, c.user_id, c.cat_breed_id, c.created_at 
		FROM cats c
		JOIN cat_house_occupancy cho ON c.id = cho.cat_id
		WHERE cho.cat_house_id = ?
		ORDER BY cho.occupied_at DESC`

	rows, err := r.db.Query(query, houseID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var cats []models.Cat
	for rows.Next() {
		var cat models.Cat
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Age, &cat.Description, &cat.UserID, &cat.CatBreedID, &cat.CreatedAt)
		if err != nil {
			return nil, nil, err
		}
		cats = append(cats, cat)
	}

	return house, cats, nil
}
