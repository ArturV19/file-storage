package config

import "os"

const (
	defaultHTTPServHost = ""
	defaultHTTPServPort = "8080"
)

const (
	envNameHTTPServHost = "HTTP_HOST"
	envNameHTTPServPort = "HTTP_PORT"
)

type HTTPConfig struct {
	Host string
	Port string
}

func newDefaultHTTPConfig() HTTPConfig {
	return HTTPConfig{
		Host: defaultHTTPServHost,
		Port: defaultHTTPServPort,
	}
}

func (c *HTTPConfig) parseEnv() {
	envServHost := os.Getenv(envNameHTTPServHost)
	if envServHost != "" {
		c.Host = envServHost
	}

	envServPort := os.Getenv(envNameHTTPServPort)
	if envServPort != "" {
		c.Port = envServPort
	}
}
