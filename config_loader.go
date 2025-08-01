package goconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// ConfigLoader loads and caches YAML configuration into a hashmap
type ConfigLoader struct {
	config      map[string]interface{}
	configPaths []string
	configName  string
	configType  string
	defaults    map[string]interface{}
}

// NewConfigLoader creates a new ConfigLoader instance
func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{
		config:      make(map[string]interface{}),
		configPaths: []string{},
		configName:  "application",
		configType:  "yaml",
		defaults:    make(map[string]interface{}),
	}
}

// SetDefault sets a default value for a key
func (c *ConfigLoader) SetDefault(key string, value interface{}) {
	c.defaults[strings.ToLower(key)] = value
}

// SetConfigName sets the name of the config file without extension
func (c *ConfigLoader) SetConfigName(name string) {
	c.configName = name
}

// SetConfigType sets the type of the config file
func (c *ConfigLoader) SetConfigType(configType string) {
	c.configType = configType
}

// AddConfigPath adds a path to search for the config file in
func (c *ConfigLoader) AddConfigPath(path string) {
	c.configPaths = append(c.configPaths, path)
}

// ReadInConfig reads in the config file from one of the search paths
func (c *ConfigLoader) ReadInConfig() error {
	// Apply defaults first
	for k, v := range c.defaults {
		c.config[k] = v // Keys in defaults are already lowercase
	}

	var configFile string
	var found bool

	for _, path := range c.configPaths {
		possibleFile := filepath.Join(path, fmt.Sprintf("%s.%s", c.configName, c.configType))
		if _, err := os.Stat(possibleFile); err == nil {
			configFile = possibleFile
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("config file not found in paths: %v", c.configPaths)
	}

	// Read the file
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	// Initialize temporary config map
	yamlConfig := make(map[string]interface{})

	// Unmarshal the YAML into the config map
	if err := yaml.Unmarshal(yamlFile, &yamlConfig); err != nil {
		return fmt.Errorf("error unmarshaling config: %v", err)
	}

	// Convert keys to strings and lowercase
	processedConfig := convertKeysToLowerStrings(yamlConfig)

	// Merge with the config map (keeping the multi-level structure)
	for k, v := range processedConfig {
		c.config[k] = v
	}

	return nil
}

func (c *ConfigLoader) GetValue(key string) (interface{}, bool) {
	key = strings.ToLower(key)
	val, exists := c.config[key]
	return val, exists
}

// convertKeysToLowerStrings recursively converts map keys to lowercase strings
func convertKeysToLowerStrings(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range m {
		// Convert the key to lowercase
		lowerKey := strings.ToLower(k)

		switch val := v.(type) {
		case map[interface{}]interface{}:
			// Convert keys to lowercase strings for this level
			strMap := make(map[string]interface{})
			for ik, iv := range val {
				// Convert interface{} key to lowercase string
				strKey := strings.ToLower(fmt.Sprintf("%v", ik))
				strMap[strKey] = iv
			}
			// Recursively convert deeper levels
			result[lowerKey] = convertKeysToLowerStrings(strMap)
		case map[string]interface{}:
			// Already has string keys, convert to lowercase and process nested maps
			result[lowerKey] = convertKeysToLowerStrings(val)
		default:
			// Not a map, keep as is
			result[lowerKey] = v
		}
	}

	return result
}
