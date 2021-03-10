package entities

import (
	"time"
)

// UserGender int
type UserGender int

// UserGender Master
const (
	Male UserGender = iota + 1
	Female
)

// UserRole string
type UserRole string

// UserRole Master
const (
	Admin    UserRole = "ADMIN"
	Customer UserRole = "CUSTOMER"
)

// Users struct
type Users struct {
	ID          int        `db:"id" json:"id"`
	Email       string     `db:"email" json:"email"`
	PhoneNumber string     `db:"phone_number" json:"phone_number"`
	FullName    string     `db:"full_name" json:"full_name"`
	Gender      UserGender `db:"gender" json:"gender"`
	Role        UserRole   `db:"role" json:"role"`
	Password    string     `db:"password" json:"password"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

// UsersLog struct
type UsersLog struct {
	ID          int        `db:"id" json:"id"`
	UserID      int        `db:"user_id" json:"user_id"`
	Email       string     `db:"email" json:"email"`
	PhoneNumber string     `db:"phone_number" json:"phone_number"`
	FullName    string     `db:"full_name" json:"full_name"`
	Gender      UserGender `db:"gender" json:"gender"`
	Role        UserRole   `db:"role" json:"role"`
	Password    string     `db:"password" json:"password"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

// Profile struct
type Profile struct {
	ID          int        `db:"id" json:"id"`
	Email       string     `db:"email" json:"email"`
	PhoneNumber string     `db:"phone_number" json:"phone_number"`
	FullName    string     `db:"full_name" json:"full_name"`
	Gender      UserGender `db:"gender" json:"gender"`
	Role        UserRole   `db:"role" json:"role"`
}
