package goconfig

import (
	"github.com/newrelic/go-agent"
)

type ConfigReader interface {
	GetValue(string) string
	GetOptionalValue(key string, defaultValue string) string
	GetIntValue(string) int
	GetOptionalIntValue(key string, defaultValue int) int
	GetFeature(key string) bool
}

type ConfigManager interface {
	ConfigReader
	Load() error
	LoadWithOptions(options map[string]interface{}) error
}

type BaseConfig struct {
	loader      *ConfigLoader
	nrConfig    newrelic.Config
	dbConfig    *DBConfig
	hasNewRelic bool
	hasDB       bool
}

func NewBaseConfig() ConfigManager {
	return &BaseConfig{
		loader:      NewConfigLoader(),
		hasNewRelic: false,
		hasDB:       false,
	}
}

func (cfg *BaseConfig) Load() error {
	return cfg.LoadWithOptions(map[string]interface{}{})
}

func (cfg *BaseConfig) LoadWithOptions(options map[string]interface{}) error {
	cfg.loader = NewConfigLoader()
	cfg.loader.SetDefault("port", "3000")
	cfg.loader.SetDefault("log_level", "warn")
	cfg.loader.SetDefault("redis_password", "")
	cfg.loader.SetConfigName("application")
	if options["configPath"] != nil {
		cfg.loader.AddConfigPath(options["configPath"].(string))
	} else {
		cfg.loader.AddConfigPath("./")
		cfg.loader.AddConfigPath("../")
	}
	cfg.loader.SetConfigType("yml")
	err := cfg.loader.ReadYamlConfig()
	if err != nil {
		return err
	}

	if options["newrelic"] != nil && options["newrelic"].(bool) {
		cfg.nrConfig = getNewRelicConfigOrPanic(cfg.loader)
		cfg.hasNewRelic = true
	}
	if options["db"] != nil && options["db"].(bool) {
		cfg.dbConfig = LoadDbConf(cfg.loader)
		cfg.hasDB = true
	}
	return nil
}

func (cfg *BaseConfig) setTestDBUrl(dbConf *DBConfig) {
	dbConf.url = getStringOrPanic(cfg.loader, "db_url_test")
	dbConf.slaveUrl = getStringOrPanic(cfg.loader, "db_url_test")
}

func (cfg *BaseConfig) LoadTestConfig(options map[string]interface{}) error {
	err := cfg.LoadWithOptions(options)
	if err != nil {
		return err
	}
	if options["db"] != nil && options["db"].(bool) && cfg.hasDB {
		cfg.setTestDBUrl(cfg.dbConfig)
	}
	return nil
}

func (cfg *BaseConfig) Newrelic() newrelic.Config {
	return cfg.nrConfig
}

func (cfg *BaseConfig) DBConfig() *DBConfig {
	return cfg.dbConfig
}

func (cfg *BaseConfig) GetValue(key string) string {
	return getStringOrPanic(cfg.loader, key)
}

func (cfg *BaseConfig) GetOptionalValue(key string, defaultValue string) string {
	v, ok := cfg.loader.GetValue(key)
	if !ok {
		return defaultValue
	}
	return v.(string)
}

func (cfg *BaseConfig) GetIntValue(key string) int {
	return getIntOrPanic(cfg.loader, key)
}

func (cfg *BaseConfig) GetOptionalIntValue(key string, defaultValue int) int {
	v, ok := cfg.loader.GetValue(key)
	if !ok {
		return defaultValue
	}
	return v.(int)
}

func (cfg *BaseConfig) GetFeature(key string) bool {
	return getFeature(cfg.loader, key)
}
