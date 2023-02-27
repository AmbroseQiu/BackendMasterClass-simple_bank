package repository

import (
	"context"

	"github.com/backendmaster/simple_bank/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresqlSessionRepository struct {
	db *gorm.DB
}

func NewpostgresqlSessionRepository(db *gorm.DB) domain.SessionRepository {
	return &postgresqlSessionRepository{
		db: db,
	}
}

func (p *postgresqlSessionRepository) Create(ctx context.Context, session domain.Session) (domain.Session, error) {
	result := p.db.Create(&session)
	if result.Error != nil {
		return domain.Session{}, result.Error
	}
	return session, nil
}

func (p *postgresqlSessionRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Session, error) {
	session := domain.Session{}
	result := p.db.First(&session, id)
	if result.Error != nil {
		return domain.Session{}, result.Error
	}
	return session, nil
}
