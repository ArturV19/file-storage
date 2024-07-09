package storage

import (
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"log"
)

var ErrUserAlreadyExists = errors.New("user already exists")

func (s *Storage) CreateUser(ctx context.Context, login, password string) (int64, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return -1, err
	}
	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil && errRollback != pgx.ErrTxClosed {
				log.Printf("storage.CreateUser tx.Rollback error: %v", errRollback)
			}
		} else {
			if errCommit := tx.Commit(ctx); errCommit != nil {
				log.Printf("storage.CreateUser tx.Commit error: %v", errCommit)
			}
		}
	}()

	passwordHash := hashPassword(password)

	var newUserID int64
	err = tx.QueryRow(ctx, "INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id", login, passwordHash).Scan(&newUserID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return -1, ErrUserAlreadyExists
		}
		return -1, err
	}

	return newUserID, nil
}
