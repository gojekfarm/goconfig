package goconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type ConfigLoader struct {
	config      map[string]interface{}
	configPaths []string
	configName  string
}

func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{
		config:      make(map[string]interface{}),
		configPaths: []string{},
		configName:  "application",
	}
}

func (c *ConfigLoader) SetDefault(key string, value interface{}) {
	c.config[strings.ToLower(key)] = value
}

func (c *ConfigLoader) SetConfigName(name string) {
	c.configName = name
}

func (c *ConfigLoader) AddConfigPath(path string) {
	c.configPaths = append(c.configPaths, path)
}

func (c *ConfigLoader) ReadYamlConfig() error {
	configFile, found := c.getConfigFile()

	if !found {
		return fmt.Errorf("config file not found in paths: %v", c.configPaths)
	}

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	yamlConfig := make(map[string]interface{})

	if err := yaml.Unmarshal(yamlFile, &yamlConfig); err != nil {
		return fmt.Errorf("error unmarshaling config: %v", err)
	}

	processedConfig := convertKeysToLowerStrings(yamlConfig)

	for k, v := range processedConfig {
		c.config[k] = v
	}

	return nil
}

func (c *ConfigLoader) getConfigFile() (string, bool) {
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

func (c *ConfigLoader) GetValue(key string) (interface{}, bool) {
	key = strings.ToLower(key)
	envVal, ok := os.LookupEnv(key)
	if ok {
		return envVal, true
	}

	val, exists := c.config[key]
	return val, exists
}

func convertKeysToLowerStrings(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range m {
		lowerKey := strings.ToLower(k)
		result[lowerKey] = v
	}

	return result
}
