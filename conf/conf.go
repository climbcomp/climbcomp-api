package conf

import (
	"log"

	"github.com/spf13/viper"
)

var (
	config *Config
)

type Config struct {
	Address     string `mapstructure:"address"`
	DatabaseUrl string `mapstructure:"database_url"`
	LogFormat   string `mapstructure:"log_format"`
	LogLevel    string `mapstructure:"log_level"`
}

// Config constructor
func NewConfig() *Config {
	viper.AddConfigPath("/etc/climbcomp")
	viper.SetConfigName("config")
	viper.SetEnvPrefix("climbcomp")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("Unable to read config: %s \n", err)
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		log.Panicf("Unable to unmarshal config: %s \n", err)
	}

	return config
}

// Returns the Config singleton instance
func Instance() *Config {
	if config == nil {
		config = NewConfig()
	}
	return config
}
