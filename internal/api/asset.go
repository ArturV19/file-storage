package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ArturV19/file-storage/internal/dto"
	"github.com/ArturV19/file-storage/internal/storage"
)

func (a *API) assetHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := a.authorize(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	assetName := strings.TrimPrefix(r.URL.Path, "/api/asset/")
	if assetName == "" {
		http.Error(w, "Name is empty", http.StatusBadRequest)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Asset not found: %s", assetName),
		})
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
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		if errors.Is(err, storage.ErrAssetNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Asset not found: %s", assetName),
			})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Internal server error",
			})
		}
		return
	}

	resp := dto.DeleteAssetResponse{
		Status: "ok",
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal server error",
		})
	}
}
