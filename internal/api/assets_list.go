package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ArturV19/file-storage/internal/dto"
)

const (
	limitDefault  = 10
	offsetDefault = 0
)

func (a *API) getUserAssetsListHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := a.authorize(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	limit, offset, err := parsePaginationParams(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	files, err := a.assetStorage.GetUserAssetsList(r.Context(), userID, limit, offset)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp := dto.GetUserAssetsListResponse{
		Assets: files,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func parsePaginationParams(r *http.Request) (int, int, error) {
	limit := limitDefault
	offset := offsetDefault

	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l < 0 {
			return 0, 0, errors.New("invalid limit parameter")
		}
		limit = l
	}

	offsetStr := r.URL.Query().Get("offset")
	if offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err != nil || o < 0 {
			return 0, 0, errors.New("invalid offset parameter")
		}
		offset = o
	}

	return limit, offset, nil
}
