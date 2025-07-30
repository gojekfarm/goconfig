package goconfig

import (
	"github.com/newrelic/go-agent"
	"github.com/spf13/viper"
	"sync"
)

type Config interface {
	GetValue(string) string
	GetOptionalValue(key string, defaultValue string) string
	GetIntValue(string) int
	GetOptionalIntValue(key string, defaultValue int) int
	GetFeature(key string) bool
}

type configuration map[string]interface{}

type BaseConfig struct {
	config      configuration
	configMutex sync.RWMutex
}

func NewBaseConfig() Config {
	return &BaseConfig{
		config: make(configuration),
	}
}

func (self *BaseConfig) Load() {
	self.LoadWithOptions(map[string]interface{}{})
}

func (self *BaseConfig) LoadWithOptions(options map[string]interface{}) {
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
	viper.ReadInConfig()
	self.config = configuration{}
	if options["newrelic"] != nil && options["newrelic"].(bool) {
		self.config["newrelic"] = getNewRelicConfigOrPanic()
	}
	if options["db"] != nil && options["db"].(bool) {
		self.config["db_config"] = LoadDbConf()
	}
}

func (self *BaseConfig) setTestDBUrl(dbConf *DBConfig) {
	dbConf.url = getStringOrPanic("db_url_test")
	dbConf.slaveUrl = getStringOrPanic("db_url_test")
}

func (self *BaseConfig) LoadTestConfig(options map[string]interface{}) error {
	self.LoadWithOptions(options)
	if options["db"] != nil && options["db"].(bool) {
		self.setTestDBUrl(self.config["db_config"].(*DBConfig))
	}
	return nil
}

func (self *BaseConfig) Newrelic() newrelic.Config {
	return self.config["newrelic"].(newrelic.Config)
}

func (self *BaseConfig) DBConfig() *DBConfig {
	return self.config["db_config"].(*DBConfig)
}

func (self *BaseConfig) GetValue(key string) string {
	self.configMutex.RLock()
	v, ok := self.config[key]
	if !ok {
		self.configMutex.RUnlock()

		v = getStringOrPanic(key)

		self.configMutex.Lock()
		self.config[key] = v
		self.configMutex.Unlock()

		return v.(string)
	}
	self.configMutex.RUnlock()
	return v.(string)
}

func (self *BaseConfig) GetOptionalValue(key string, defaultValue string) string {
	self.configMutex.RLock()
	v, ok := self.config[key]
	if !ok {
		self.configMutex.RUnlock()

		if v = viper.GetString(key); !viper.IsSet(key) {
			v = defaultValue
		}

		self.configMutex.Lock()
		self.config[key] = v
		self.configMutex.Unlock()

		return v.(string)
	}
	self.configMutex.RUnlock()
	return v.(string)
}

func (self *BaseConfig) GetIntValue(key string) int {
	self.configMutex.RLock()
	v, ok := self.config[key]
	if !ok {
		self.configMutex.RUnlock()
		v = getIntOrPanic(key)

		self.configMutex.Lock()
		self.config[key] = v
		self.configMutex.Unlock()

		return v.(int)
	}
	self.configMutex.RUnlock()
	return v.(int)
}

func (self *BaseConfig) GetOptionalIntValue(key string, defaultValue int) int {
	self.configMutex.RLock()
	v, ok := self.config[key]
	if !ok {
		self.configMutex.RUnlock()

		if v = viper.GetInt(key); !viper.IsSet(key) {
			v = defaultValue
		}

		self.configMutex.Lock()
		self.config[key] = v
		self.configMutex.Unlock()

		return v.(int)
	}
	self.configMutex.RUnlock()
	return v.(int)
}

func (self *BaseConfig) GetFeature(key string) bool {
	self.configMutex.RLock()
	v, ok := self.config[key]
	if !ok {
		self.configMutex.RUnlock()

		v = getFeature(key)

		self.configMutex.Lock()
		self.config[key] = v
		self.configMutex.Unlock()

		return v.(bool)
	}
	self.configMutex.RUnlock()
	return v.(bool)
}
