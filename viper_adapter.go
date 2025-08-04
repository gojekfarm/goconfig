package goconfig

// ViperAccessor Interface that defines used and available viper interactions in this package.
// Functions in this interface is implicitly implemented by the viper package.
type ViperAccessor interface {
	SetDefault(key string, value interface{})
	Get(key string) interface{}
	GetString(key string) string
	GetInt(key string) int
	Set(key string, value interface{})
	ReadInConfig() error
	AutomaticEnv()
	SetConfigName(in string)
	SetConfigType(extension string)
	AddConfigPath(path string)
}

type ConfigAccessorToViperAdapter struct {
	viper ViperAccessor
}

func (c ConfigAccessorToViperAdapter) Get(key string) (interface{}, bool) {
	result := c.viper.Get(key)
	if result != nil {
		return result, true
	}
	return nil, false
}

func (c ConfigAccessorToViperAdapter) Set(key string, value interface{}) {
	c.viper.Set(key, value)
}

func (c ConfigAccessorToViperAdapter) Load() error {
	return c.viper.ReadInConfig()
}

func (c ConfigAccessorToViperAdapter) SetConfigName(name string) {
	c.viper.SetConfigName(name)
}

// AddExtension sets the config type for the viper instance.
// This is a workaround since Viper does not support multiple extensions directly.
// It allows setting the config type to a single extension, which can be used
// to load configurations with that specific extension.
func (c ConfigAccessorToViperAdapter) AddExtension(extension string) {
	c.viper.SetConfigType(extension)
}

func (c ConfigAccessorToViperAdapter) AddPath(path string) {
	c.viper.AddConfigPath(path)
}

func NewConfigAccessorToViperAdapter(viper ViperAccessor) ConfigFileAccessor {
	return &ConfigAccessorToViperAdapter{
		viper: viper,
	}
}

type ViperConfigFileAccessorAdapter struct {
	configFileAccessor ConfigFileAccessor
}

func (cfa *ViperConfigFileAccessorAdapter) SetConfigType(extension string) {
	cfa.configFileAccessor.AddExtension(extension)
}

func (cfa *ViperConfigFileAccessorAdapter) AddConfigPath(path string) {
	cfa.configFileAccessor.AddPath(path)
}

func (cfa *ViperConfigFileAccessorAdapter) ReadInConfig() error {
	return cfa.configFileAccessor.Load()
}

func (cfa *ViperConfigFileAccessorAdapter) AutomaticEnv() {
	// pass, no automatic env handling in ConfigFileAccessor
	return
}

func (cfa *ViperConfigFileAccessorAdapter) SetConfigName(in string) {
	cfa.configFileAccessor.SetConfigName(in)
}

func (cfa *ViperConfigFileAccessorAdapter) SetDefault(key string, value interface{}) {
	cfa.configFileAccessor.Set(key, value)
}

func (cfa *ViperConfigFileAccessorAdapter) Get(key string) interface{} {
	value, _ := cfa.configFileAccessor.Get(key)
	return value
}

func (cfa *ViperConfigFileAccessorAdapter) GetString(key string) string {
	value, _ := cfa.configFileAccessor.Get(key)
	return value.(string)
}

func (cfa *ViperConfigFileAccessorAdapter) GetInt(key string) int {
	value, _ := cfa.configFileAccessor.Get(key)
	return value.(int)
}

func (cfa *ViperConfigFileAccessorAdapter) Set(key string, value interface{}) {
	cfa.configFileAccessor.Set(key, value)
}

func NewViperConfigFileAccessorAdapter(configFileAccessor ConfigFileAccessor) ViperAccessor {
	return &ViperConfigFileAccessorAdapter{
		configFileAccessor: configFileAccessor,
	}
}
