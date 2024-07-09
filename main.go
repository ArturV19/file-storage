package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ArturV19/file-storage/config"
	"github.com/ArturV19/file-storage/internal/api"
	"github.com/ArturV19/file-storage/internal/storage"
	"github.com/ArturV19/file-storage/internal/utils"
)

func main() {
	err := utils.LoadEnv(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg := config.NewDefaultConfig()
	if err := cfg.ParseEnv(); err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("[" + strings.ToUpper(cfg.LogLvl) + "] ")

	newStorage, err := storage.New(cfg.Storage)
	if err != nil {
		log.Fatalf("Error initializing storage: %v\n", err)
	} else {
		log.Println("AssetStorage initialized")
	}

	// Start the expired session removal routine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go newStorage.RemoveExpiredSessions(ctx)

	apiServer := api.New(cfg.HTTPConfig, newStorage, newStorage)
	if err != nil {
		log.Fatalf("Error initializing API server: %v\n", err)
	}

	shutdownSig := make(chan os.Signal, 1)
	signal.Notify(shutdownSig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctxStart, ctxStartCancel := context.WithCancel(context.Background())

	errServingCh := make(chan error)
	go func() {
		errServing := apiServer.Start(ctxStart)
		errServingCh <- errServing
	}()

	select {
	case shutdownSigValue := <-shutdownSig:
		close(shutdownSig)
		log.Printf("Shutdown signal received: %s", strings.ToUpper(shutdownSigValue.String()))
	case errServing := <-errServingCh:
		if errServing != nil {
			log.Printf("Error from API server: %v", errServing)
		}
	}

	ctxStartCancel()

	ctxClose, ctxCloseCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCloseCancel()

	if err = apiServer.GracefulStop(ctxClose); err != nil {
		log.Printf("Error during graceful stop: %v", err)
		if err == context.DeadlineExceeded {
			return
		}
	} else {
		log.Println("API server gracefully stopped")
	}

	if err = newStorage.Close(ctxClose); err != nil {
		log.Printf("Error closing storage: %v", err)
		if err == context.DeadlineExceeded {
			return
		}
	} else {
		log.Println("AssetStorage closed")
	}
}
