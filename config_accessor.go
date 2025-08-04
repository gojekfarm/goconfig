package goconfig

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

type ConfigAccessor struct {
	repository  ConfigRepository
	configPaths []string
	configName  string
}

func NewConfigAccessor() *ConfigAccessor {
	return &ConfigAccessor{
		repository:  NewInMemoryConfigRepository(),
		configPaths: []string{},
		configName:  "application",
	}
}

func (c *ConfigAccessor) Set(key string, value interface{}) {
	c.repository.Set(key, value)
}

func (c *ConfigAccessor) SetConfigName(name string) {
	c.configName = name
}

func (c *ConfigAccessor) AddPath(path string) {
	c.configPaths = append(c.configPaths, path)
}

func (c *ConfigAccessor) Load() error {
	configFile, found := c.getConfigFile()

	if !found {
		return fmt.Errorf("config file not found in paths: %v", c.configPaths)
	}

	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	yamlConfig := make(map[string]interface{})

	if err := yaml.Unmarshal(yamlFile, &yamlConfig); err != nil {
		return fmt.Errorf("error unmarshaling config: %v", err)
	}

	for k, v := range yamlConfig {
		c.repository.Set(k, v)
	}

	return nil
}

func (c *ConfigAccessor) getConfigFile() (string, bool) {
	extensions := []string{"yml", "yaml"}
	for _, path := range c.configPaths {
		for _, ext := range extensions {
			possibleFile := filepath.Join(path, fmt.Sprintf("%s.%s", c.configName, ext))
			if _, err := os.Stat(possibleFile); err == nil {
				return possibleFile, true
			}
		}
	}
	return "", false
}

func (c *ConfigAccessor) Get(key string) (interface{}, bool) {
	envVal, ok := os.LookupEnv(key)
	if ok {
		return envVal, true
	}
	key = toLowercaseKey(key)
	lcaseEnvVal, ok := os.LookupEnv(key)
	if ok {
		return lcaseEnvVal, true
	}

	val, exists := c.repository.Get(key)
	return val, exists
}
