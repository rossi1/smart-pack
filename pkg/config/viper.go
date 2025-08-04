package config

import (
	"strings"

	"github.com/rossi1/smart-pack/config"
	"github.com/spf13/viper"
)

// ViperConfig takes care of loading application name, configuration.
type ViperConfig struct {
}

// LoadConfig takes care of loading application name, configuration.
func (c *ViperConfig) LoadConfig(path string, cfg config.Config) error {
	for k, v := range cfg.Defaults() {
		viper.SetDefault(k, v)
	}
	viper.AddConfigPath(path)
	viper.SetConfigName(cfg.Name())
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.Unmarshal(cfg)
}
