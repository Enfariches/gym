package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

type Storage struct {
	db *gorm.DB
}

func New(databaseURL string) (*Storage, error) {
	const op = "postgres.New"

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) SaveAdmin(ctx context.Context, email string, passHash []byte) error {
	const op = "postgres.SaveAdmin"
	
	

	return nil
}
