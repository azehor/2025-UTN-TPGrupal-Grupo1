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

	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "development"
	}
	envLower := strings.ToLower(env)

	// Solo cargamos .env en desarrollo
	if envLower == "development" {
		_ = godotenv.Load() // ignoramos error, usamos variables del sistema si no existe
	}

	host := os.Getenv("APP_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	portStr := os.Getenv("APP_PORT")

	if portStr == "" {
		if envLower == "production" {
			return nil, fmt.Errorf("APP_PORT is required in production")

		}
		// default en development.
		return &Config{
			Env:  envLower,
			Host: host,
			Port: 8080}, nil
	}

	// Parsear puerto
	p, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid APP_PORT: %v", err)
	}

	return &Config{
		Env:  envLower,
		Host: host,
		Port: p,
	}, nil
}
