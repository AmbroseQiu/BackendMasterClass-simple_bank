package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/backendmaster/simple_bank/domain"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type postgresqlUserRepository struct {
	db *gorm.DB
}

func NewpostgresqlUserRepository(db *gorm.DB) domain.UsersRepository {
	return &postgresqlUserRepository{
		db: db,
	}
}

func (p *postgresqlUserRepository) GetByUsername(cxt context.Context, username string) (*domain.User, error) {
	user := &domain.User{
		Username: username,
	}
	// check user is existed and check password
	result := p.db.First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) || errors.Is(result.Error, sql.ErrNoRows) {
			return nil, domain.ErrorUserNotFound
		}
		return nil, domain.ErrorInternalServerError
	}
	return user, nil
}

func (p *postgresqlUserRepository) Create(cxt context.Context, user domain.User) (*domain.User, error) {

	// repository.Create(user)
	result := p.db.Create(&user)
	if result.Error != nil {
		if pqErr, ok := result.Error.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, domain.ErrorUniqueViolation
			}
		}
		return nil, domain.ErrorInternalServerError
	}
	return &user, nil
}

func (p *postgresqlUserRepository) Update(cxt context.Context, user domain.User) (*domain.User, error) {
	result := p.db.Save(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// func (p *postgresqlUserRepository) PrintLog() string {
// 	return "HI"
// }
