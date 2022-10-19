package msnet

import "fmt"

type ClientConfig struct {
}

type ConfigProtocol interface {
	GetBaseUrl() string
	GetReadTimeout() int
	GetConnectTimeout() int
	GetApis() map[string]Api
	FullPath(api Api) string
}

type Config struct {
	BaseUrl        string         `mapstructure:"baseUrl"`
	ReadTimeout    int            `mapstructure:"readTimeout"`
	ConnectTimeout int            `mapstructure:"connectTimeout"`
	Apis           map[string]Api `mapstructure:"apis"`
}

func (config *Config) GetBaseUrl() string {
	return config.BaseUrl
}

func (config *Config) GetReadTimeout() int {
	return config.ReadTimeout
}

func (config *Config) GetConnectTimeout() int {
	return config.ConnectTimeout
}

func (config *Config) GetApis() map[string]Api {
	return config.Apis
}

func (config *Config) FullPath(api Api) string {
	return fmt.Sprintf("%v%v", config.BaseUrl, api.Uri)
}

type Api struct {
	Uri    string `mapstructure:"uri"`
	Method string `mapstructure:"method"`
}
