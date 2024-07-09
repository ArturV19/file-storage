package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/ArturV19/file-storage/internal/storage"
	"github.com/ArturV19/file-storage/internal/types"
)

const MAX_UPLOAD_SIZE = 1 * 1024 * 1024 * 1024 // 1GB

func (a *API) UploadAssetHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	assetName := strings.TrimPrefix(r.URL.Path, "/api/upload-asset/")
	if assetName == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userID, err := a.userStorage.ValidateToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing request body: %v", err)
		}
	}(r.Body)

	if r.ContentLength > MAX_UPLOAD_SIZE {
		http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
		return
	}

	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
		return
	}

	var originalName string
	var body io.Reader

	if strings.HasPrefix(mediaType, "application/json") {
		originalName = ""
		body = r.Body
	} else {
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Invalid file", http.StatusBadRequest)
			return
		}
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				http.Error(w, "File close error", http.StatusInternalServerError)
				return
			}
		}(file)

		originalName = header.Filename
		body = file
	}

	uploadAssetData := types.UploadAssetData{
		UserID:       userID,
		AssetName:    assetName,
		OriginalName: originalName,
		ContentType:  mediaType,
		Body:         body,
	}

	err = a.assetStorage.UploadAsset(r.Context(), uploadAssetData)
	if err != nil {
		if errors.Is(err, storage.ErrAssetAlreadyExists) {
			http.Error(w, "Asset already exists", http.StatusConflict)
		} else if errors.Is(err, io.ErrUnexpectedEOF) {
			http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		return
	}
}
