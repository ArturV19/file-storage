package storage

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/ArturV19/file-storage/config"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(cfg config.StorageConfig) (*Storage, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	log.Printf("Connecting to database with connection string: %s\n", connStr)

	db, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	log.Println("Successfully connected to the database")

	return &Storage{db: db}, nil
}

func (s *Storage) Close(_ context.Context) error {
	s.db.Close()
	return nil
}

func hashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}
