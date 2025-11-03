package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config centraliza todas las variables de configuración
type Config struct {
	Env  string
	Host string
	Port int
}

// LoadConfig carga la configuración desde variables de entorno o .env en desarrollo
func LoadConfig() (*Config, error) {
	// Detectar entorno: APP_ENV (preferida) o GO_ENV; default: development
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = os.Getenv("GO_ENV")
	}
	if env == "" {
		env = "development"
	}
	envLower := strings.ToLower(env)

	// Solo cargamos .env en desarrollo/local
	if envLower == "development" || envLower == "dev" || envLower == "local" {
		_ = godotenv.Load() // ignoramos error, usamos variables del sistema si no existe
	}

	// Host
	host := os.Getenv("APP_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	// Port (soportar APP_PORT y PORT)
	portStr := os.Getenv("APP_PORT")
	if portStr == "" {
		portStr = os.Getenv("PORT")
	}
	port := 8080 // default
	if portStr != "" {
		p, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid APP_PORT: %v", err)
		}
		port = p
	} else if envLower == "production" || envLower == "prod" {
		return nil, fmt.Errorf("APP_PORT is required in production")
	}

	return &Config{
		Env:  envLower,
		Host: host,
		Port: port,
	}, nil
}
