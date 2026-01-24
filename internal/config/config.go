package config

import (
	"os"
	"strings"
)

type Config struct {
	GeminiAPIKey      string
	GeminiAPIEndpoint string
	Port              string
	AllowedOrigins    []string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	origins := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	if origins != "" {
		allowedOrigins = strings.Split(origins, ",")
	} else {
		allowedOrigins = []string{"*"}
	}

	endpoint := os.Getenv("GEMINI_API_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://127.0.0.1:8045"
	}

	return &Config{
		GeminiAPIKey:      os.Getenv("GEMINI_API_KEY"),
		GeminiAPIEndpoint: endpoint,
		Port:              port,
		AllowedOrigins:    allowedOrigins,
	}
}
