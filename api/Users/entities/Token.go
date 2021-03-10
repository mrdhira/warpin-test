package entities

import "github.com/dgrijalva/jwt-go"

// TokenClaim struct
type TokenClaim struct {
	UserID   int      `json:"user_id"`
	UserRole UserRole `json:"user_role"`
	jwt.StandardClaims
}
