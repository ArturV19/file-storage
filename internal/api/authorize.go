package api

import (
	"errors"
	"net/http"
	"strings"
)

func (a *API) authorize(r *http.Request) (int64, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, errors.New("missing authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return 0, errors.New("missing token")
	}

	userID, err := a.userStorage.ValidateToken(r.Context(), token)
	if err != nil {
		return 0, errors.New("invalid token")
	}

	return userID, nil
}
