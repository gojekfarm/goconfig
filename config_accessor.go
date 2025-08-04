package goconfig

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

type ConfigAccessor interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
}

type ConfigFileAccessor interface {
	ConfigAccessor
	Load() error
	SetConfigName(name string)
	AddExtension(extension string)
	AddPath(path string)
}

type YamlConfigAccessor struct {
	levelDelimiter string
	repository     ConfigRepository
	configPaths    []string
	configName     string
	configFileExts []string
}

func NewYamlConfigAccessor(
	repository ConfigRepository,
	configPaths []string,
	configName string,
	configFileExts []string,
	multiLevelDelimiter string,
) ConfigFileAccessor {
	return &YamlConfigAccessor{
		repository:     repository,
		configPaths:    configPaths,
		configName:     configName,
		configFileExts: configFileExts,
		levelDelimiter: multiLevelDelimiter,
	}
}

// NewDefaultYamlConfigAccessor creates a new YamlConfigAccessor with default settings.
// It initializes the repository with an in-memory config repository decorated with environment variable support.
func NewDefaultYamlConfigAccessor() ConfigFileAccessor {
	return &YamlConfigAccessor{
		repository: NewEnvConfigRepositoryDecorator(
			NewInMemoryConfigRepository(),
		),
		levelDelimiter: ".",
		configPaths:    []string{},
		configName:     "application",
		configFileExts: []string{"yaml", "yml"},
	}
}

func (c *YamlConfigAccessor) SetConfigName(name string) {
	c.configName = name
}

func (c *YamlConfigAccessor) AddExtension(extension string) {
	c.configFileExts = append(c.configFileExts, extension)
}

func (c *YamlConfigAccessor) AddPath(path string) {
	c.configPaths = append(c.configPaths, path)
}

func (c *YamlConfigAccessor) Load() error {
	if c.configPaths == nil || len(c.configPaths) == 0 {
		return fmt.Errorf("no config paths set, please add paths using AddPath method")
	}

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

	for k, v := range c.flattenMap(yamlConfig) {
		c.repository.Set(k, v)
	}

	return nil
}

func (c *YamlConfigAccessor) getConfigFile() (string, bool) {
	for _, path := range c.configPaths {
		for _, ext := range c.configFileExts {
			possibleFile := filepath.Join(path, fmt.Sprintf("%s.%s", c.configName, ext))
			if _, err := os.Stat(possibleFile); err == nil {
				return possibleFile, true
			}
		}
	}
	return "", false
}

func (c *YamlConfigAccessor) flattenMap(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	var flatten func(map[string]interface{}, string)
	flatten = func(nested map[string]interface{}, prefix string) {
		for k, v := range nested {
			key := k
			if prefix != "" {
				key = prefix + c.levelDelimiter + k
			}

			switch val := v.(type) {
			case map[string]interface{}:
				flatten(val, key)
			case map[interface{}]interface{}:
				stringMap := make(map[string]interface{})
				for mk, mv := range val {
					if mkString, ok := mk.(string); ok {
						stringMap[mkString] = mv
					} else {
						fmt.Printf("skipping non-string key during flatten map: %v\n", mk)
					}
				}
				flatten(stringMap, key)
			default:
				result[key] = v
			}
		}
	}

	flatten(m, "")
	return result
}

func (c *YamlConfigAccessor) Get(key string) (interface{}, bool) {
	return c.repository.Get(key)
}

func (c *YamlConfigAccessor) Set(key string, value interface{}) {
	c.repository.Set(key, value)
}

func (c *YamlConfigAccessor) SetDefault(key string, value interface{}) {
	c.repository.Set(key, value)
}
