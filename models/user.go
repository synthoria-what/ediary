package models

import "time"

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type GetUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
	Role     string `json:"role"`
}

type DummyUser struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

type UpdateUser struct {
	Username     *string `json:"username"`
	PasswordHash *string
	Email        *string `json:"email"`
	Role         *string `json:"role"`
}

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserCreateForm struct {
	Username string
	Password string
	Role     string
}

type UserStats struct {
	RequestCount int
	IsBanned     bool
}
