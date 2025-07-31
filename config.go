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

type configuration map[string]interface{}

type BaseConfig struct {
	config      configuration
	configMutex sync.RWMutex
}

func NewBaseConfig() ConfigManager {
	return &BaseConfig{
		config: make(configuration),
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
	cfg.config = configuration{}
	if options["newrelic"] != nil && options["newrelic"].(bool) {
		cfg.config["newrelic"] = getNewRelicConfigOrPanic()
	}
	if options["db"] != nil && options["db"].(bool) {
		cfg.config["db_config"] = LoadDbConf()
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
		cfg.setTestDBUrl(cfg.config["db_config"].(*DBConfig))
	}
	return nil
}

func (cfg *BaseConfig) Newrelic() newrelic.Config {
	return cfg.config["newrelic"].(newrelic.Config)
}

func (cfg *BaseConfig) DBConfig() *DBConfig {
	return cfg.config["db_config"].(*DBConfig)
}

func (cfg *BaseConfig) GetValue(key string) string {
	cfg.configMutex.RLock()
	v, ok := cfg.config[key]
	cfg.configMutex.RUnlock()

	if !ok {
		cfg.configMutex.Lock()
		v, ok = cfg.config[key]
		if !ok {
			v = getStringOrPanic(key)
			cfg.config[key] = v
		}
		cfg.configMutex.Unlock()

		return v.(string)
	}
	return v.(string)
}

func (cfg *BaseConfig) GetOptionalValue(key string, defaultValue string) string {
	cfg.configMutex.RLock()
	v, ok := cfg.config[key]
	cfg.configMutex.RUnlock()

	if !ok {
		cfg.configMutex.Lock()
		v, ok := cfg.config[key]
		if !ok {
			if v = viper.GetString(key); !viper.IsSet(key) {
				v = defaultValue
			}
			cfg.config[key] = v
		}
		cfg.configMutex.Unlock()

		return v.(string)
	}
	return v.(string)
}

func (cfg *BaseConfig) GetIntValue(key string) int {
	cfg.configMutex.RLock()
	v, ok := cfg.config[key]
	cfg.configMutex.RUnlock()

	if !ok {
		cfg.configMutex.Lock()
		v, ok = cfg.config[key]
		if !ok {
			v = getIntOrPanic(key)
			cfg.config[key] = v
		}
		cfg.configMutex.Unlock()

		return v.(int)
	}
	return v.(int)
}

func (cfg *BaseConfig) GetOptionalIntValue(key string, defaultValue int) int {
	cfg.configMutex.RLock()
	v, ok := cfg.config[key]
	cfg.configMutex.RUnlock()

	if !ok {
		cfg.configMutex.Lock()
		v, ok := cfg.config[key]
		if !ok {
			if v = viper.GetInt(key); !viper.IsSet(key) {
				v = defaultValue
			}
			cfg.config[key] = v
		}
		cfg.configMutex.Unlock()

		return v.(int)
	}
	return v.(int)
}

func (cfg *BaseConfig) GetFeature(key string) bool {
	cfg.configMutex.RLock()
	v, ok := cfg.config[key]
	cfg.configMutex.RUnlock()

	if !ok {
		cfg.configMutex.Lock()
		v, ok := cfg.config[key]
		if !ok {
			v = getFeature(key)
			cfg.config[key] = v
		}
		cfg.configMutex.Unlock()

		return v.(bool)
	}
	return v.(bool)
}
