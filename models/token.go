package models

import "github.com/golang-jwt/jwt/v5"

type TokenClaims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`

	jwt.RegisteredClaims
}

type LoginResponse struct {
	Token string `json:"token"`
}
