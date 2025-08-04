package goconfig

import (
	"github.com/newrelic/go-agent"
	"github.com/spf13/viper"
	"sync"
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
