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

type InMemoryConfigRepository struct {
	config map[string]interface{}
}

func NewInMemoryConfigRepository() ConfigRepository {
	return &InMemoryConfigRepository{
		config: make(map[string]interface{}),
	}
}

func (repo *InMemoryConfigRepository) Get(key string) (interface{}, bool) {
	value, exists := repo.config[toLowercaseKey(key)]
	return value, exists
}

func (repo *InMemoryConfigRepository) Set(key string, value interface{}) {
	repo.config[toLowercaseKey(key)] = value
}

func toLowercaseKey(key string) string {
	return strings.ToLower(key)
}
