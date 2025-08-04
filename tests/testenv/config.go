package testenv

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rossi1/smart-pack/config"
	"github.com/spf13/viper"
)

var (
	ErrEnvFileNotFound = errors.New("environment file not found")
	ErrChangeDirectory = errors.New("error while changing base directory")
)

type Config struct {
	Cfg       *config.AppConfig
	PubSubURI string
}

type Option func(*Config)

// envFile is referenced from the root directory.
const envFile = ".test.ci.env"

func NewConfig(options ...Option) (*Config, error) {
	found, err := findAndSetWorkingDirectory(envFile)
	if err != nil {
		return nil, fmt.Errorf("change directory: %w", err)
	}

	if !found {
		return nil, ErrEnvFileNotFound
	}

	v := viper.New()
	v.SetConfigFile(envFile)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg config.AppConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	appConfig := &Config{
		Cfg: &cfg,
	}

	for _, opt := range options {
		opt(appConfig)
	}

	return appConfig, nil
}

func WithDatabaseURL(url string) Option {
	return func(c *Config) {
		c.Cfg.DatabaseURL = url
	}
}

func findAndSetWorkingDirectory(filename string) (bool, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return false, err
	}

	maxDepth := 3
	for i := 0; i <= maxDepth; i++ {
		fullPath := filepath.Join(cwd, filename)
		if _, err := os.Stat(fullPath); err == nil {
			if err := os.Chdir(cwd); err != nil {
				return false, err
			}
			return true, nil
		}
		cwd = filepath.Dir(cwd)
	}
	return false, nil
}
