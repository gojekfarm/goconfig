package goconfig

import (
	"fmt"
	"github.com/spf13/cast"
)

func getIntOrPanic(accessor ConfigAccessor, key string) int {
	v, ok := accessor.Get(key)
	if !ok {
		panic(fmt.Errorf("%s key is not set", key))
	}
	return cast.ToInt(v)
}

func getFeature(accessor ConfigAccessor, key string) bool {
	v, ok := accessor.Get(key)
	if !ok {
		return false
	}
	return cast.ToBool(v)
}

func getStringOrPanic(accessor ConfigAccessor, key string) string {
	v, ok := accessor.Get(key)
	if !ok {
		panic(fmt.Errorf("%s key is not set", key))
	}
	return cast.ToString(v)
}

func getString(accessor ConfigAccessor, key string) string {
	v, _ := accessor.Get(key)
	return cast.ToString(v)
}
