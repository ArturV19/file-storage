package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"time"

	"github.com/ArturV19/file-storage/internal/types"
)

const SESSION_DURATION = 24 * time.Hour

var ErrInvalidLoginPassword = errors.New("invalid login/password")
var ErrInvalidToken = errors.New("invalid token")

func (s *Storage) NewSession(ctx context.Context, authData types.AuthData) (string, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil && errRollback != pgx.ErrTxClosed {
				fmt.Printf("storage.NewSession tx.Rollback error: %v", errRollback)
			}
		}
	}()

	var userID int64
	err = tx.QueryRow(ctx, "SELECT id FROM users WHERE login = $1 AND password_hash = $2", authData.Login, hashPassword(authData.Password)).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrInvalidLoginPassword
		}
		return "", err
	}

	var token string
	err = tx.QueryRow(ctx, "SELECT id FROM sessions WHERE uid = $1 AND expires_at > now()", userID).Scan(&token)
	if err == nil {
		if errCommit := tx.Commit(ctx); errCommit != nil {
			return "", errCommit
		}
		return token, nil
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", err
	}

	_, err = tx.Exec(ctx, "DELETE FROM sessions WHERE uid = $1", userID)
	if err != nil {
		return "", err
	}

	token, err = generateToken()
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(SESSION_DURATION)

	_, err = tx.Exec(ctx, "INSERT INTO sessions (id, uid, ip_address, expires_at) VALUES ($1, $2, $3, $4)", token, userID, authData.IPAddress, expiresAt)
	if err != nil {
		return "", err
	}

	if err = tx.Commit(ctx); err != nil {
		return "", err
	}

	return token, nil
}

func (s *Storage) ValidateToken(ctx context.Context, token string) (int64, error) {
	var userID int64
	err := s.db.QueryRow(ctx, "SELECT uid FROM sessions WHERE id = $1 AND expires_at > now()", token).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrInvalidToken
		}
		return 0, err
	}
	return userID, nil
}

func (s *Storage) RemoveExpiredSessions(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_, err := s.db.Exec(ctx, "DELETE FROM sessions WHERE expires_at <= now()")
			if err != nil {
				log.Printf("Failed to remove expired sessions: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func generateToken() (string, error) {
	bytes := make([]byte, 16) // для 128-битного токена
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
