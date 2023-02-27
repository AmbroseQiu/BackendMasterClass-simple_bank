package usercase

import (
	"context"

	"github.com/backendmaster/simple_bank/domain"
	"github.com/google/uuid"
)

type sessionUseCase struct {
	sessionrepo domain.SessionRepository
}

func NewSessionUseCase(sessionrepo domain.SessionRepository) domain.SessionUseCase {
	return &sessionUseCase{
		sessionrepo: sessionrepo,
	}
}

func (s *sessionUseCase) CreateSession(ctx context.Context, session domain.Session) (domain.Session, error) {
	session, err := s.sessionrepo.Create(ctx, session)
	if err != nil {
		return domain.Session{}, err
	}
	return session, nil
}

func (s *sessionUseCase) GetSessionByID(ctx context.Context, id uuid.UUID) (domain.Session, error) {
	session, err := s.sessionrepo.GetByID(ctx, id)
	if err != nil {
		return domain.Session{}, err
	}
	return session, nil
}
