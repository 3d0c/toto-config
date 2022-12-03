package config

import (
	"sync"
	"time"
)

var (
	instance *Config
	once     sync.Once
)

// Server config secon
type Server struct {
	APIVersion  string
	Address     string
	ReadTimeout time.Duration
	WriteTemout time.Duration
}

// Logger config section
type Logger struct {
	Level       string
	AddCaller   bool
	OutputPaths []string
}

// Database config section
type Database struct {
	DSN string
}

// Config is a complete config structure
type Config struct {
	Server
	Logger
	Database
}

// TheConfig config singleton
func TheConfig() *Config {
	once.Do(func() {
		instance = new(Config)
	})

	return instance
}
