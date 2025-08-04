package goconfig_test

import (
	"github.com/gojekfarm/goconfig/v2"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type AppConfig struct {
	goconfig.BaseConfig
}

func TestShouldSetDefaultForPort(t *testing.T) {
	config := &AppConfig{}
	config.Load()
	assert.Equal(t, "3000", config.GetValue("port"))
}

func TestShouldSetDefaultForLogLevel(t *testing.T) {
	config := &AppConfig{}
	config.Load()
	assert.Equal(t, "warn", config.GetValue("log_level"))
}

func TestShouldSetDefaultForRedisPassword(t *testing.T) {
	config := &AppConfig{}
	config.Load()
	assert.Equal(t, "", config.GetValue("redis_password"))
}

func TestShouldSetNewRelicBasedOnApplicationConfig(t *testing.T) {
	config := &AppConfig{}
	config.LoadWithOptions(map[string]interface{}{"newrelic": true})
	assert.Equal(t, "foo", config.Newrelic().AppName)
	assert.Equal(t, "bar", config.Newrelic().License)
}

func TestShouldGetValueBasedOnApplicationConfig(t *testing.T) {
	config := &AppConfig{}
	config.Load()
	assert.Equal(t, "bar", config.GetValue("foo"))
}

func TestShouldGetOptionalValueBasedOnApplicationConfig(t *testing.T) {
	config := &AppConfig{}
	config.Load()
	assert.Equal(t, "bar", config.GetOptionalValue("foo", "baz"))
	assert.Equal(t, "bar", config.GetOptionalValue("foo", "baz"))
	assert.Equal(t, "baz", config.GetOptionalValue("bar", "baz"))
}

func TestShouldGetIntValueBasedOnApplicationConfig(t *testing.T) {
	config := &AppConfig{}
	config.Load()
	assert.Equal(t, 1, config.GetIntValue("someInt"))
	assert.Equal(t, 1, config.GetIntValue("someInt"))
}

func TestShouldGetOptionalIntValueBasedOnApplicationConfig(t *testing.T) {
	config := &AppConfig{}
	config.Load()
	assert.Equal(t, 1, config.GetOptionalIntValue("someInt", 10))
	assert.Equal(t, 1, config.GetOptionalIntValue("someInt", 10))
	assert.Equal(t, 10, config.GetOptionalIntValue("bar", 10))
}

func TestShouldGetFeature(t *testing.T) {
	config := &AppConfig{}
	config.Load()
	assert.True(t, config.GetFeature("someFeature"))
	assert.True(t, config.GetFeature("someFeature"))
	assert.False(t, config.GetFeature("someOtherFeature"))
	assert.False(t, config.GetFeature("someUnknownFeature"))
}

func TestShouldGetDBConfig(t *testing.T) {
	config := &AppConfig{}
	config.LoadWithOptions(map[string]interface{}{"db": true})
	assert.Equal(t, "test://something", config.DBConfig().Url())
	assert.Equal(t, "test://something", config.DBConfig().SlaveUrl())
	assert.Equal(t, "postgres", config.DBConfig().Driver())
	assert.Equal(t, 5, config.DBConfig().MaxConn())
	assert.Equal(t, 2, config.DBConfig().IdleConn())
	assert.Equal(t, time.Duration(1000000000), config.DBConfig().ConnMaxLifetime())
}

func TestShouldGetTestDBConfigOnLoadTestConfig(t *testing.T) {
	config := &AppConfig{}
	config.LoadTestConfig(map[string]interface{}{"db": true})
	assert.Equal(t, "test://somethingTest", config.DBConfig().Url())
	assert.Equal(t, "test://somethingTest", config.DBConfig().SlaveUrl())
	assert.Equal(t, "postgres", config.DBConfig().Driver())
	assert.Equal(t, 5, config.DBConfig().MaxConn())
	assert.Equal(t, 2, config.DBConfig().IdleConn())
	assert.Equal(t, time.Duration(1000000000), config.DBConfig().ConnMaxLifetime())
}

func TestShouldSetConfigPathBasedOnOptionaParam(t *testing.T) {
	config := &AppConfig{}
	confData := []byte("foo: 9998\n")
	ioutil.WriteFile("/tmp/application.yml", confData, 0644)
	config.LoadWithOptions(map[string]interface{}{"configPath": "/tmp"})
	assert.Equal(t, "9998", config.GetValue("foo"))
}

func TestShouldNestedConfigValueBasedOnApplicationConfig(t *testing.T) {
	config := &AppConfig{}
	config.Load()
	assert.Equal(t, "value", config.GetValue("nested_config.level_1.level_2.level_3"))
	assert.Equal(t, 42, config.GetIntValue("nested_config.level_1.level_2.level_4"))
	assert.Equal(t, true, config.GetFeature("nested_config.level_1.level_5"))
	assert.Equal(t, false, config.GetFeature("nested_config.level_6"))
}
