package config

import "time"

type Loader interface {
	LoadConfig(path string, cfg Config) error
}

type Config interface {
	Name() string
	Defaults() map[string]string
}

type AppConfig struct {
	Hostname               string        `mapstructure:"HOSTNAME"`
	ApplicationAPITimeout  time.Duration `mapstructure:"APPLICATION_API_TIMEOUT"`
	ApplicationName        string        `mapstructure:"APPLICATION_NAME"`
	ApplicationEnvironment string        `mapstructure:"APPLICATION_ENV"`
	LogLevel               string        `mapstructure:"LOG_LEVEL"`
	CORSAllowedOrigins     string        `mapstructure:"CORS_ALLOWED_ORIGINS"`
	Port                   int           `mapstructure:"PORT"`
	DatabaseURL            string        `mapstructure:"DATABASE_URL"`
	DatabaseMigrationPath  string        `mapstructure:"DATABASE_MIGRATION_PATH"`
}

func (c *AppConfig) Name() string {
	return "app"
}

func (c *AppConfig) Defaults() map[string]string {
	return map[string]string{
		"HOSTNAME":                "localhost",
		"APPLICATION_API_TIMEOUT": "30s",
		"APPLICATION_NAME":        "smart-pack",
		"APPLICATION_ENV":         "development",
		"LOG_LEVEL":               "info",
		"CORS_ALLOWED_ORIGINS":    "*",
		"PORT":                    "8080",
		"DATABASE_URL":            "postgres://smartpack:smartpack@localhost:5432/smartpack?sslmode=disable",
		"DATABASE_MIGRATION_PATH": "resources/db",
	}
}
