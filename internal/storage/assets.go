package storage

import (
	"context"
	"errors"
	"io"

	"github.com/ArturV19/file-storage/internal/types"
)

var ErrAssetAlreadyExists = errors.New("asset already exists")
var ErrAssetNotFound = errors.New("asset not found")

func (s *Storage) UploadAsset(ctx context.Context, uploadAssetData types.UploadAssetData) error {
	data, err := io.ReadAll(uploadAssetData.Body)
	if err != nil {
		return err
	}

	var exists bool
	err = s.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM assets WHERE name = $1 AND uid = $2)", uploadAssetData.AssetName, uploadAssetData.UserID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return ErrAssetAlreadyExists
	}

	_, err = s.db.Exec(ctx, "INSERT INTO assets (name, original_name, uid, data, content_type) VALUES ($1, $2, $3, $4, $5)", uploadAssetData.AssetName, uploadAssetData.OriginalName, uploadAssetData.UserID, data, uploadAssetData.ContentType)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetAsset(ctx context.Context, userID int64, assetName string) ([]byte, string, string, error) {
	var data []byte
	var originalName, contentType string
	err := s.db.QueryRow(ctx, "SELECT data, original_name, content_type FROM assets WHERE name = $1 AND uid = $2", assetName, userID).Scan(&data, &originalName, &contentType)

	if err != nil {
		return nil, "", "", err
	}

	return data, originalName, contentType, nil
}

func (s *Storage) GetUserAssetsList(ctx context.Context, userID int64, limit, offset int) ([]types.Asset, error) {
	rows, err := s.db.Query(ctx, "SELECT name, original_name, content_type, created_at FROM assets WHERE uid = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3", userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []types.Asset

	for rows.Next() {
		var file types.Asset
		err := rows.Scan(&file.Name, &file.OriginalName, &file.ContentType, &file.CreatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

func (s *Storage) DeleteAsset(ctx context.Context, userID int64, assetName string) error {
	var exists bool
	err := s.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM assets WHERE name = $1 AND uid = $2)", assetName, userID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrAssetNotFound
	}

	_, err = s.db.Exec(ctx, "DELETE FROM assets WHERE name = $1 AND uid = $2", assetName, userID)
	if err != nil {
		return err
	}

	return nil
}
