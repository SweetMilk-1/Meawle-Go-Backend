package services

import (
	"errors"
	"time"

	"meawle/internal/models"
	"meawle/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists        = errors.New("email already exists")
	ErrUnauthorized       = errors.New("unauthorized")
)

// JWTClaims представляет claims для JWT токена
type JWTClaims struct {
	UserID  int    `json:"user_id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// UserService представляет сервис для работы с пользователями
type UserService struct {
	repo      repositories.UserRepository
	jwtSecret string
}

// NewUserService создает новый экземпляр сервиса пользователей
func NewUserService(repo repositories.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

// Register регистрирует нового пользователя
func (s *UserService) Register(req *models.UserCreateRequest) (*models.UserResponse, error) {
	// Проверяем существование email
	exists, err := s.repo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailExists
	}

	// Создаем пользователя
	user := &models.User{
		Email:    req.Email,
		Password: req.Password, // В реальном приложении здесь должно быть хеширование пароля
		IsAdmin:  req.IsAdmin,
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// Login выполняет вход пользователя и возвращает JWT токен
func (s *UserService) Login(req *models.UserLoginRequest) (string, *models.UserResponse, error) {
	// Получаем пользователя по email

	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	// Проверяем пароль (в реальном приложении должно быть сравнение хешей)
	if user.Password != req.Password {
		return "", nil, ErrInvalidCredentials
	}

	// Генерируем JWT токен
	token, err := s.generateJWT(user)
	if err != nil {
		return "", nil, err
	}

	response := user.ToResponse()
	return token, &response, nil
}

// GetUserByID возвращает пользователя по ID
func (s *UserService) GetUserByID(id int) (*models.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	response := user.ToResponse()
	return &response, nil
}

// GetAllUsers возвращает всех пользователей
func (s *UserService) GetAllUsers() ([]models.UserResponse, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []models.UserResponse
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}

	return responses, nil
}

// UpdateUser обновляет данные пользователя
func (s *UserService) UpdateUser(id int, req *models.UserUpdateRequest) error {
	// Проверяем существование пользователя
	_, err := s.repo.GetByID(id)
	if err != nil {
		return ErrUserNotFound
	}

	// Если обновляется email, проверяем его уникальность
	if req.Email != nil {
		exists, err := s.repo.ExistsByEmail(*req.Email)
		if err != nil {
			return err
		}
		if exists {
			return ErrEmailExists
		}
	}

	return s.repo.Update(id, req)
}

// DeleteUser удаляет пользователя
func (s *UserService) DeleteUser(id int) error {
	// Проверяем существование пользователя
	_, err := s.repo.GetByID(id)
	if err != nil {
		return ErrUserNotFound
	}

	return s.repo.Delete(id)
}

// ValidateToken проверяет JWT токен и возвращает claims
func (s *UserService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, ErrUnauthorized
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrUnauthorized
}

// generateJWT генерирует JWT токен для пользователя
func (s *UserService) generateJWT(user *models.User) (string, error) {
	claims := JWTClaims{
		UserID:  user.ID,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
