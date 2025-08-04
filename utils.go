package goconfig

import (
	"fmt"
	"github.com/spf13/cast"
)

func getIntOrPanic(loader *ConfigAccessor, key string) int {
	v, ok := loader.Get(key)
	if !ok {
		panic(fmt.Errorf("%s key is not set", key))
	}
	return cast.ToInt(v)
}

func getFeature(loader *ConfigAccessor, key string) bool {
	v, ok := loader.Get(key)
	if !ok {
		return false
	}
	return cast.ToBool(v)
}

func getStringOrPanic(loader *ConfigAccessor, key string) string {
	v, ok := loader.Get(key)
	if !ok {
		panic(fmt.Errorf("%s key is not set", key))
	}
	return cast.ToString(v)
}
