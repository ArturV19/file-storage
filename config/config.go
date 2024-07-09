package config

import (
	"fmt"
	"os"
)

const (
	defaultLogLvl = "info"
	envNameLogLvl = "LOG_LVL"
)

type Config struct {
	HTTPConfig HTTPConfig
	Storage    StorageConfig
	LogLvl     string
}

func NewDefaultConfig() Config {
	return Config{
		HTTPConfig: newDefaultHTTPConfig(),
		Storage:    newDefaultStorageConfig(),
		LogLvl:     defaultLogLvl,
	}
}

func (c *Config) ParseEnv() error {
	c.HTTPConfig.parseEnv()
	c.Storage.parseEnv()

	if err := c.parseEnvLogLvl(); err != nil {
		return err
	}

	return nil
}

func (c *Config) parseEnvLogLvl() error {
	envLogLvl := os.Getenv(envNameLogLvl)
	if envLogLvl != "" {
		logLevel, err := parseLogLevel(envLogLvl)
		if err != nil {
			return fmt.Errorf("parse log lvl: %s", envLogLvl)
		}
		c.LogLvl = logLevel
	}
	return nil
}

func parseLogLevel(level string) (string, error) {
	switch level {
	case "debug", "info", "warn", "error", "fatal", "panic":
		return level, nil
	default:
		return "info", fmt.Errorf("unknown log level: %s", level)
	}
}
