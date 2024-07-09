package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func (a *API) GetUserAssetsListHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	userID, err := a.userStorage.ValidateToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	limit, offset := parsePaginationParams(r)

	files, err := a.assetStorage.GetUserAssetsList(r.Context(), userID, limit, offset)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(files)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func parsePaginationParams(r *http.Request) (int, int) {
	limit := 10 // default limit
	offset := 0 // default offset

	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	offsetStr := r.URL.Query().Get("offset")
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	return limit, offset
}
