package gorm

import (
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(conn *sql.DB) (*gorm.DB, error) {
	dialector := postgres.New(postgres.Config{Conn: conn})
	return gorm.Open(dialector, &gorm.Config{})
}
