package goconfig

import "strings"

type ConfigRepository interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
}

type InMemoryConfigRepositoryImpl struct {
	config map[string]interface{}
}

func NewInMemoryConfigRepository() ConfigRepository {
	return &InMemoryConfigRepositoryImpl{
		config: make(map[string]interface{}),
	}
}

func (repo *InMemoryConfigRepositoryImpl) Get(key string) (interface{}, bool) {
	value, exists := repo.config[key]
	return value, exists
}

func (repo *InMemoryConfigRepositoryImpl) Set(key string, value interface{}) {
	repo.config[toLowercaseKey(key)] = value
}

func toLowercaseKey(key string) string {
	return strings.ToLower(key)
}
