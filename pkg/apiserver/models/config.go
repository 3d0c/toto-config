package models

import (
	"database/sql"
	"errors"

	"github.com/3d0c/toto-config/pkg/database"
)

// ErrNotFound return this type to not to expose sql package into handlers
var ErrNotFound = errors.New("nothing found")

// ConfigScheme is a Config Model representation
type ConfigScheme struct {
	Package       string `json:"-"`
	CountryCode   string `json:"-"`
	PercentileMin int    `json:"-"`
	PercentileMax int    `json:"-"`
	MainSKU       string `json:"main_sku"`
}

// Config Model
type Config struct{}

// NewConfigModel is a Config Model constructor
func NewConfigModel() (*Config, error) {
	return &Config{}, nil
}

// FindBy finds SKU by package name, country code and percentile seed
func (c *Config) FindBy(packageName string, cc string, seed int) (*ConfigScheme, error) {
	var (
		cs   *ConfigScheme = new(ConfigScheme)
		err  error
		stmt string = "SELECT main_sku from toto_configuration where package = ? and country_code = ? and percentile_min > ? and percentile_max <= ?"
	)

	if err = database.Instance().QueryRow(stmt, packageName, cc, seed, seed).Scan(&cs.MainSKU); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return cs, nil
}
