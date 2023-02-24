package domain

import (
	"context"
	"time"

	"github.com/backendmaster/simple_bank/token"
	"github.com/google/uuid"
)

type User struct {
	// gorm.Model
	Username          string `gorm:"primary_key"`
	HashedPassword    string
	FullName          string
	Email             string
	PasswordChangedAt time.Time
	CreatedAt         time.Time
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,alphanumunicode"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type LoginUserRequest struct {
	Username             string `json:"username" binding:"required,alphanumunicode"`
	Password             string `json:"password" binding:"required,min=6"`
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

type LoginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResponse `json:"user_response"`
}

type UsersRepository interface {
	GetByUsername(cxt context.Context, username string) (User, error)
	Create(cxt context.Context, user User) (User, error)
	Update(cxt context.Context, user User) (User, error)
	// PrintLog() string
}

type UsersTableUseCase interface {
	CreateUser(cxt context.Context, req CreateUserRequest) (UserResponse, error)
	CreateToken(username string, duration time.Duration) (string, *token.Payload, error)
	// PrintLog() string
	LoginUser(cxt context.Context, req LoginUserRequest) (LoginUserResponse, error)
}
