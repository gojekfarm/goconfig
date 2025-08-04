package goconfig

import (
	"os"
	"strings"
)

type ConfigRepository interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
}

type EnvConfigRepositoryDecorator struct {
	innerRepo ConfigRepository
}

func NewEnvConfigRepositoryDecorator(innerRepo ConfigRepository) ConfigRepository {
	return &EnvConfigRepositoryDecorator{innerRepo}
}

func (d *EnvConfigRepositoryDecorator) Get(key string) (interface{}, bool) {
	envVal, ok := os.LookupEnv(key)
	if ok {
		return envVal, true
	}
	lcaseEnvVal, ok := os.LookupEnv(toLowercaseKey(key))
	if ok {
		return lcaseEnvVal, true
	}
	return d.innerRepo.Get(key)
}

func (d *EnvConfigRepositoryDecorator) Set(key string, value interface{}) {
	d.innerRepo.Set(key, value)
	return
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
