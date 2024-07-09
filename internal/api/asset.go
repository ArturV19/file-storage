package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/ArturV19/file-storage/internal/storage"
)

func (a *API) AssetHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	assetName := strings.TrimPrefix(r.URL.Path, "/api/asset/")
	if assetName == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userID, err := a.userStorage.ValidateToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		a.handleGetAsset(w, r, userID, assetName)
	case http.MethodDelete:
		a.handleDeleteAsset(w, r, userID, assetName)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (a *API) handleGetAsset(w http.ResponseWriter, r *http.Request, userID int64, assetName string) {
	data, originalName, contentType, err := a.assetStorage.GetAsset(r.Context(), userID, assetName)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", contentType)
	if originalName != "" {
		w.Header().Set("Content-Disposition", "attachment; filename=\""+originalName+"\"")
	}
	_, err = w.Write(data)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (a *API) handleDeleteAsset(w http.ResponseWriter, r *http.Request, userID int64, assetName string) {
	err := a.assetStorage.DeleteAsset(r.Context(), userID, assetName)
	if err != nil {
		if errors.Is(err, storage.ErrAssetNotFound) {
			http.Error(w, "Asset not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
