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
	PrivateKey  string
	Certificate string
	Insecure    bool
}

// Logger config section
type Logger struct {
	Level       string
	AddCaller   bool
	OutputPaths []string
}

// Database config section
type Database struct {
	DSN     string
	Dialect string
}

// GeoIP database config section
type GeoIP struct {
	UseGeoIP   bool
	DBFileName string
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
