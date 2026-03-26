package gorm

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(db *sqlx.DB) (*gorm.DB, error) {
	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	return gormDb, nil
}
