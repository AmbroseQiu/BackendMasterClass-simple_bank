package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type SessionRepository interface {
	Create(ctx context.Context, session Session) (Session, error)
	GetByID(ctx context.Context, id uuid.UUID) (Session, error)
}

type SessionUseCase interface {
	CreateSession(ctx context.Context, session Session) (Session, error)
	GetSessionByID(ctx context.Context, id uuid.UUID) (Session, error)
}
