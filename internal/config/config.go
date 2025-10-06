package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Treasury TreasuryConfig `yaml:"treasury"`
}

type TreasuryConfig struct {
	GrantAmountPerUser int    `yaml:"grant_amount_per_user"`
	Username           string `yaml:"username"`
	GrantDelayMs       int    `yaml:"grant_delay_ms"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
}

type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(config); err != nil {
		return nil, err
	}

	if dsn := os.Getenv("DATABASE_DSN"); dsn != "" {
		config.Database.DSN = dsn
	}

	return config, nil
}
