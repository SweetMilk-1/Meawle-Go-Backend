package models

// User представляет модель пользователя
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"` // Пароль не должен сериализоваться в JSON
	IsAdmin  bool   `json:"is_admin"`
}

// UserCreateRequest представляет данные для создания пользователя
type UserCreateRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	IsAdmin  bool   `json:"is_admin"`
}

// UserUpdateRequest представляет данные для обновления пользователя
type UserUpdateRequest struct {
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=6"`
	IsAdmin  *bool   `json:"is_admin,omitempty"`
}

// UserLoginRequest представляет данные для входа пользователя
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse представляет ответ с данными пользователя
type UserResponse struct {
	ID      int    `json:"id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

// ToResponse преобразует User в UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:      u.ID,
		Email:   u.Email,
		IsAdmin: u.IsAdmin,
	}
}
