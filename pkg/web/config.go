package web

import (
	"cryptobackup/pkg/storage"
)

// ServerConfig holds the configuration for the web server
type ServerConfig struct {
	Host        string           // Server host (default: 0.0.0.0)
	Port        int              // Server port (default: 8080)
	StoragePath string           // Path to storage directory
	Username    string           // Login username
	PasswordHash string          // Hashed login password
	Storage     storage.Storage  // Storage instance
	SessionSecret string         // Secret for session management
}

// NewServerConfig creates a new server configuration with defaults
func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Host: "0.0.0.0",
		Port: 8080,
		StoragePath: "./backup",
	}
}
