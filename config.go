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
	Accessor ConfigFileAccessor
	nrConfig newrelic.Config
	dbConfig *DBConfig
}

func NewBaseConfig(accessor ConfigFileAccessor) ConfigManager {
	return &BaseConfig{
		Accessor: accessor,
	}
}

func (cfg *BaseConfig) Load() error {
	return cfg.LoadWithOptions(map[string]interface{}{})
}

func (cfg *BaseConfig) LoadWithOptions(options map[string]interface{}) error {
	if cfg.Accessor == nil {
		cfg.Accessor = NewDefaultYamlConfigAccessor()
	}
	cfg.Accessor.Set("port", "3000")
	cfg.Accessor.Set("log_level", "warn")
	cfg.Accessor.Set("redis_password", "")
	cfg.Accessor.SetConfigName("application")
	if options["configPath"] != nil {
		cfg.Accessor.AddPath(options["configPath"].(string))
	} else {
		cfg.Accessor.AddPath("./")
		cfg.Accessor.AddPath("../")
	}
	err := cfg.Accessor.Load()
	if err != nil {
		return err
	}

	if options["newrelic"] != nil && options["newrelic"].(bool) {
		cfg.nrConfig = getNewRelicConfigOrPanic(cfg.Accessor)
	}
	if options["db"] != nil && options["db"].(bool) {
		cfg.dbConfig = LoadDbConf(cfg.Accessor)
	}
	return nil
}

func (cfg *BaseConfig) setTestDBUrl(dbConf *DBConfig) {
	dbConf.url = getStringOrPanic(cfg.Accessor, "db_url_test")
	dbConf.slaveUrl = getStringOrPanic(cfg.Accessor, "db_url_test")
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
	return getStringOrPanic(cfg.Accessor, key)
}

func (cfg *BaseConfig) GetOptionalValue(key string, defaultValue string) string {
	v, ok := cfg.Accessor.Get(key)
	if !ok {
		return defaultValue
	}
	return v.(string)
}

func (cfg *BaseConfig) GetIntValue(key string) int {
	return getIntOrPanic(cfg.Accessor, key)
}

func (cfg *BaseConfig) GetOptionalIntValue(key string, defaultValue int) int {
	v, ok := cfg.Accessor.Get(key)
	if !ok {
		return defaultValue
	}
	return v.(int)
}

func (cfg *BaseConfig) GetFeature(key string) bool {
	return getFeature(cfg.Accessor, key)
}
