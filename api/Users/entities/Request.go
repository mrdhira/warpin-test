package entities

// RegisterRequest struct
type RegisterRequest struct {
	Email       string     `json:"email" validate:"required"`
	PhoneNumber string     `json:"phone_number" validate:"-"`
	FullName    string     `json:"full_name" validate:"required"`
	Gender      UserGender `json:"gender" validate:"required"`
	Role        UserRole   `json:"role" validate:"required"`
	Password    string     `json:"password" validate:"required"`
}

// LoginRequest struct
type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// ProfileRequest struct
type ProfileRequest struct {
	UserID int `json:"user_id" validate:"required"`
}

// UpdateProfileRequest struct
type UpdateProfileRequest struct {
	UserID      int    `json:"user_id" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"-"`
	FullName    string `json:"full_name" validate:"-"`
}

// UpdatePasswordRequest struct
type UpdatePasswordRequest struct {
	UserID   int    `json:"user_id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UsersByIDRequest struct
type UsersByIDRequest struct {
	UserID int `json:"user_id" validate:"required"`
}
