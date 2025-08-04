package goconfig

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

type YamlConfigAccessor struct {
	repository  ConfigRepository
	configPaths []string
	configName  string
}

func NewYamlConfigAccessor() *YamlConfigAccessor {
	return &YamlConfigAccessor{
		repository: NewEnvConfigRepositoryDecorator(
			NewInMemoryConfigRepository(),
		),
		configPaths: []string{},
		configName:  "application",
	}
}

func (c *YamlConfigAccessor) Set(key string, value interface{}) {
	c.repository.Set(key, value)
}

func (c *YamlConfigAccessor) SetConfigName(name string) {
	c.configName = name
}

func (c *YamlConfigAccessor) AddPath(path string) {
	c.configPaths = append(c.configPaths, path)
}

func (c *YamlConfigAccessor) Load() error {
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

func (c *YamlConfigAccessor) getConfigFile() (string, bool) {
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

func (c *YamlConfigAccessor) Get(key string) (interface{}, bool) {
	val, exists := c.repository.Get(key)
	return val, exists
}
