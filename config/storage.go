package config

import (
	"os"
)

const (
	defaultStorageHost       = "127.0.0.1"
	defaultStoragePort       = "5432"
	defaultStorageUser       = "admin"
	defaultStoragePassword   = "password"
	defaultFileStorageDBName = "file_storage"
)
const (
	envNameStorageHost     = "DB_HOST"
	envNameStoragePort     = "DB_PORT"
	envNameStorageUser     = "DB_USER"
	envNameStoragePassword = "DB_PASSWORD"
	envNameDBName          = "DB_NAME"
)

type StorageConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func newDefaultStorageConfig() StorageConfig {
	return StorageConfig{
		Host:     defaultStorageHost,
		Port:     defaultStoragePort,
		User:     defaultStorageUser,
		Password: defaultStoragePassword,
		DBName:   defaultFileStorageDBName,
	}
}

func (c *StorageConfig) parseEnv() {
	envHost := os.Getenv(envNameStorageHost)
	if envHost != "" {
		c.Host = envHost
	}

	envPort := os.Getenv(envNameStoragePort)
	if envPort != "" {
		c.Port = envPort
	}

	envUser := os.Getenv(envNameStorageUser)
	if envUser != "" {
		c.User = envUser
	}

	envPassword := os.Getenv(envNameStoragePassword)
	if envPassword != "" {
		c.Password = envPassword
	}

	envDBName := os.Getenv(envNameDBName)
	if envDBName != "" {
		c.DBName = envDBName
	}
}
