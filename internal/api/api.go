package api

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/ArturV19/file-storage/config"
	"github.com/ArturV19/file-storage/internal/types"
)

type assetStorage interface {
	UploadAsset(ctx context.Context, uploadAssetData types.UploadAssetData) error
	GetAsset(ctx context.Context, userID int64, assetName string) ([]byte, string, string, error)
	DeleteAsset(ctx context.Context, userID int64, assetName string) error
	GetUserAssetsList(ctx context.Context, userID int64, limit, offset int) ([]types.Asset, error)
}

type userStorage interface {
	CreateUser(ctx context.Context, login, password string) (int64, error)
	ValidateToken(ctx context.Context, token string) (int64, error)
	NewSession(ctx context.Context, authData types.AuthData) (string, error)
}

type API struct {
	httpServer   *http.Server
	assetStorage assetStorage
	userStorage  userStorage
}

func New(httpCfg config.HTTPConfig, assetStorage assetStorage, userStorage userStorage) (newAPI *API) {
	newAPI = &API{
		assetStorage: assetStorage,
		userStorage:  userStorage,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/users/create", newAPI.createUserHandler)
	mux.HandleFunc("/api/auth", newAPI.authenticateUserHandler)
	mux.HandleFunc("/api/upload-asset/", newAPI.UploadAssetHandler)
	mux.HandleFunc("/api/asset/", newAPI.AssetHandler)
	mux.HandleFunc("/api/list-assets", newAPI.GetUserAssetsListHandler)

	newAPI.httpServer = &http.Server{
		Addr:    net.JoinHostPort(httpCfg.Host, httpCfg.Port),
		Handler: mux,
	}

	return newAPI
}

func (a *API) Start(ctx context.Context) error {
	httpServeEndSig := make(chan struct{})
	go func() {
		log.Printf("Starting HTTP server on %s\n", a.httpServer.Addr)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v\n", err)
		}
		close(httpServeEndSig)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-httpServeEndSig:
		return nil
	}
}

func (a *API) GracefulStop(ctx context.Context) error {
	stopCh := make(chan struct{})
	go func() {
		if err := a.httpServer.Shutdown(ctx); err != nil {
			log.Printf("HTTP server shutdown error: %v\n", err)
		}
		close(stopCh)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-stopCh:
		return nil
	}
}
