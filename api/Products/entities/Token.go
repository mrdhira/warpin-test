package entities

import "github.com/dgrijalva/jwt-go"

// UserRole string
type UserRole string

// UserRole Master
const (
	Admin    UserRole = "ADMIN"
	Customer UserRole = "CUSTOMER"
)

// TokenClaim struct
type TokenClaim struct {
	UserID   int      `json:"user_id"`
	UserRole UserRole `json:"user_role"`
	jwt.StandardClaims
}
