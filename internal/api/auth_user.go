package api

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/ArturV19/file-storage/internal/storage"
	"github.com/ArturV19/file-storage/internal/types"
)

func (a *API) authenticateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

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

	ipAddress, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "Invalid IP address", http.StatusInternalServerError)
		return
	}

	authData := types.AuthData{
		Login:     req.Login,
		Password:  req.Password,
		IPAddress: ipAddress,
	}

	token, err := a.userStorage.NewSession(ctx, authData)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidLoginPassword) {
			http.Error(w, "Invalid login/password", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	resp := map[string]interface{}{
		"token": token,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
