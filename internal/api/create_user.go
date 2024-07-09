package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/ArturV19/file-storage/internal/dto"
	"github.com/ArturV19/file-storage/internal/storage"
)

func (a *API) createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	req.Login = strings.TrimSpace(req.Login)
	if req.Login == "" {
		http.Error(w, "Empty login", http.StatusBadRequest)
		return
	}

	req.Password = strings.TrimSpace(req.Password)
	if req.Password == "" {
		http.Error(w, "Empty password", http.StatusBadRequest)
		return
	}

	userID, err := a.userStorage.CreateUser(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			http.Error(w, "User already exists", http.StatusConflict)
		} else {
			log.Printf("storage.CreateUser error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	resp := dto.CreateUserResponse{
		ID: userID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
