package repositories

import (
	"meawle/internal/models"
)

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll() ([]models.User, error)
	Update(id int, user *models.UserUpdateRequest) error
	Delete(id int) error
	ExistsByEmail(email string) (bool, error)
}

type userRepository struct {
	db Database
}

// NewUserRepository создает новый экземпляр репозитория пользователей
func NewUserRepository(db Database) UserRepository {
	return &userRepository{db: db}
}

// Create создает нового пользователя
func (r *userRepository) Create(user *models.User) error {
	query := `INSERT INTO users (email, password, is_admin) VALUES (?, ?, ?)`

	result, err := r.db.Execute(query, user.Email, user.Password, user.IsAdmin)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

// GetByID возвращает пользователя по ID
func (r *userRepository) GetByID(id int) (*models.User, error) {
	query := `SELECT id, email, password, is_admin FROM users WHERE id = ?`

	row := r.db.QueryRow(query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.IsAdmin)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmail возвращает пользователя по email
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password, is_admin FROM users WHERE email = ?`

	row := r.db.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.IsAdmin)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetAll возвращает всех пользователей
func (r *userRepository) GetAll() ([]models.User, error) {
	query := `SELECT id, email, password, is_admin FROM users ORDER BY id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.IsAdmin)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// Update обновляет данные пользователя
func (r *userRepository) Update(id int, updateReq *models.UserUpdateRequest) error {
	query := `UPDATE users SET `
	params := []interface{}{}

	if updateReq.Email != nil {
		query += "email = ?, "
		params = append(params, *updateReq.Email)
	}

	if updateReq.Password != nil {
		query += "password = ?, "
		params = append(params, *updateReq.Password)
	}

	if updateReq.IsAdmin != nil {
		query += "is_admin = ?, "
		params = append(params, *updateReq.IsAdmin)
	}

	// Убираем последнюю запятую и пробел
	query = query[:len(query)-2]
	query += " WHERE id = ?"
	params = append(params, id)

	_, err := r.db.Execute(query, params...)
	return err
}

// Delete удаляет пользователя
func (r *userRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Execute(query, id)
	return err
}

// ExistsByEmail проверяет существование пользователя по email
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ?`

	var count int
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
