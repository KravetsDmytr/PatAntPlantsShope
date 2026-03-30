package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
	JWT    JWTConfig    `yaml:"jwt"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}

func Load(path string) (*Config, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("помилка шляху конфіга: %w", err)
	}

	raw, err := os.ReadFile(abs)
	if err != nil {
		return nil, fmt.Errorf("помилка читання конфіга: %w", err)
	}

	var cfg Config
	if err = yaml.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("помилка парсингу конфіга: %w", err)
	}

	return &cfg, nil
}
