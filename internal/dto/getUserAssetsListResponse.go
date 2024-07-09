package dto

import "github.com/ArturV19/file-storage/internal/types"

type GetUserAssetsListResponse struct {
	Assets []types.Asset `json:"assets"`
}
