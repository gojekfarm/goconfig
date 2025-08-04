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
	accessor ConfigFileAccessor
	nrConfig newrelic.Config
	dbConfig *DBConfig
}

func NewBaseConfig() ConfigManager {
	return &BaseConfig{
		accessor: NewDefaultYamlConfigAccessor(),
	}
}

func (cfg *BaseConfig) Load() error {
	return cfg.LoadWithOptions(map[string]interface{}{})
}

func (cfg *BaseConfig) LoadWithOptions(options map[string]interface{}) error {
	cfg.accessor = NewDefaultYamlConfigAccessor()
	cfg.accessor.Set("port", "3000")
	cfg.accessor.Set("log_level", "warn")
	cfg.accessor.Set("redis_password", "")
	cfg.accessor.SetConfigName("application")
	if options["configPath"] != nil {
		cfg.accessor.AddPath(options["configPath"].(string))
	} else {
		cfg.accessor.AddPath("./")
		cfg.accessor.AddPath("../")
	}
	err := cfg.accessor.Load()
	if err != nil {
		return err
	}

	if options["newrelic"] != nil && options["newrelic"].(bool) {
		cfg.nrConfig = getNewRelicConfigOrPanic(cfg.accessor)
	}
	if options["db"] != nil && options["db"].(bool) {
		cfg.dbConfig = LoadDbConf(cfg.accessor)
	}
	return nil
}

func (cfg *BaseConfig) setTestDBUrl(dbConf *DBConfig) {
	dbConf.url = getStringOrPanic(cfg.accessor, "db_url_test")
	dbConf.slaveUrl = getStringOrPanic(cfg.accessor, "db_url_test")
}

func (cfg *BaseConfig) LoadTestConfig(options map[string]interface{}) error {
	err := cfg.LoadWithOptions(options)
	if err != nil {
		return err
	}
	if options["db"] != nil && options["db"].(bool) {
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
	return getStringOrPanic(cfg.accessor, key)
}

func (cfg *BaseConfig) GetOptionalValue(key string, defaultValue string) string {
	v, ok := cfg.accessor.Get(key)
	if !ok {
		return defaultValue
	}
	return v.(string)
}

func (cfg *BaseConfig) GetIntValue(key string) int {
	return getIntOrPanic(cfg.accessor, key)
}

func (cfg *BaseConfig) GetOptionalIntValue(key string, defaultValue int) int {
	v, ok := cfg.accessor.Get(key)
	if !ok {
		return defaultValue
	}
	return v.(int)
}

func (cfg *BaseConfig) GetFeature(key string) bool {
	return getFeature(cfg.accessor, key)
}
