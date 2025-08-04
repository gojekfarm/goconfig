package goconfig

import (
	"github.com/newrelic/go-agent"
	"github.com/spf13/viper"
	"sync"
)

// ConfigReader provides methods to retrieve configuration values in different formats.
// Values are cached after the first retrieval for better performance.
type ConfigReader interface {
	// GetValue retrieves a string configuration value for the given key.
	// Panics if the key doesn't exist in the configuration.
	GetValue(string) string

	// GetOptionalValue retrieves a string configuration value for the given key.
	// Returns the defaultValue if the key doesn't exist in the configuration.
	GetOptionalValue(key string, defaultValue string) string

	// GetIntValue retrieves an integer configuration value for the given key.
	// Panics if the key doesn't exist or if the value cannot be parsed as an integer.
	GetIntValue(string) int

	// GetOptionalIntValue retrieves an integer configuration value for the given key.
	// Returns the defaultValue if the key doesn't exist or cannot be parsed as an integer.
	GetOptionalIntValue(key string, defaultValue int) int

	// GetFeature retrieves a boolean feature flag value.
	// Returns false if the key doesn't exist or cannot be parsed as a boolean.
	GetFeature(key string) bool
}

// ConfigManager extends the ConfigReader interface with methods for loading configuration.
// It provides functionality to load configuration from files and environment variables,
// with support for custom loading options.
type ConfigManager interface {
	ConfigReader

	// Load initializes the configuration with default settings.
	// It reads from application.yaml files in current and parent directories
	// and from environment variables.
	// Returns an error if configuration loading fails.
	Load() error

	// LoadWithOptions initializes the configuration with custom options.
	// Options can include:
	//   - "configPath": string - Path to look for configuration files
	//   - "newrelic": bool - Whether to load New Relic configuration
	//   - "db": bool - Whether to load database configuration
	// Returns an error if configuration loading fails.
	LoadWithOptions(options map[string]interface{}) error
}

type BaseConfig struct {
	config sync.Map
}

func NewBaseConfig() ConfigManager {
	return &BaseConfig{
		config: sync.Map{},
	}
}

func (cfg *BaseConfig) Load() error {
	return cfg.LoadWithOptions(map[string]interface{}{})
}

func (cfg *BaseConfig) LoadWithOptions(options map[string]interface{}) error {
	viper.SetDefault("port", "3000")
	viper.SetDefault("log_level", "warn")
	viper.SetDefault("redis_password", "")
	viper.AutomaticEnv()
	viper.SetConfigName("application")
	if options["configPath"] != nil {
		viper.AddConfigPath(options["configPath"].(string))
	} else {
		viper.AddConfigPath("./")
		viper.AddConfigPath("../")
	}
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	cfg.config = sync.Map{}
	if options["newrelic"] != nil && options["newrelic"].(bool) {
		cfg.config.Store("newrelic", getNewRelicConfigOrPanic())
	}
	if options["db"] != nil && options["db"].(bool) {
		cfg.config.Store("db_config", LoadDbConf())
	}
	return nil
}

func (cfg *BaseConfig) setTestDBUrl(dbConf *DBConfig) {
	dbConf.url = getStringOrPanic("db_url_test")
	dbConf.slaveUrl = getStringOrPanic("db_url_test")
}

func (cfg *BaseConfig) LoadTestConfig(options map[string]interface{}) error {
	cfg.LoadWithOptions(options)
	if options["db"] != nil && options["db"].(bool) {
		dbUrl, _ := cfg.config.Load("db_config")
		cfg.setTestDBUrl(dbUrl.(*DBConfig))
	}
	return nil
}

func (cfg *BaseConfig) Newrelic() newrelic.Config {
	nrConfig, _ := cfg.config.Load("newrelic")
	return nrConfig.(newrelic.Config)
}

func (cfg *BaseConfig) DBConfig() *DBConfig {
	dbConfig, _ := cfg.config.Load("db_config")
	return dbConfig.(*DBConfig)
}

func (cfg *BaseConfig) GetValue(key string) string {
	v, ok := cfg.config.Load(key)

	if !ok {
		v = getStringOrPanic(key)
		cfg.config.Store(key, v)

		return v.(string)
	}
	return v.(string)
}

func (cfg *BaseConfig) GetOptionalValue(key string, defaultValue string) string {
	v, ok := cfg.config.Load(key)

	if !ok {
		if v = viper.GetString(key); !viper.IsSet(key) {
			v = defaultValue
		}
		cfg.config.Store(key, v)

		return v.(string)
	}
	return v.(string)
}

func (cfg *BaseConfig) GetIntValue(key string) int {
	v, ok := cfg.config.Load(key)

	if !ok {
		v = getIntOrPanic(key)
		cfg.config.Store(key, v)

		return v.(int)
	}
	return v.(int)
}

func (cfg *BaseConfig) GetOptionalIntValue(key string, defaultValue int) int {
	v, ok := cfg.config.Load(key)

	if !ok {
		if v = viper.GetInt(key); !viper.IsSet(key) {
			v = defaultValue
		}
		cfg.config.Store(key, v)

		return v.(int)
	}
	return v.(int)
}

func (cfg *BaseConfig) GetFeature(key string) bool {
	v, ok := cfg.config.Load(key)

	if !ok {
		v = getFeature(key)
		cfg.config.Store(key, v)

		return v.(bool)
	}
	return v.(bool)
}
