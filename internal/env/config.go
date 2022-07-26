package env

import (
	"github.com/spf13/viper"
)

type EnvConfig struct {
	HTTP_ADDRESS      string
	MAX_LIMITER       float64
	DB_DRIVER         string
	DATABASE_NAME     string
	DATABASE_USERNAME string
	DATABASE_PASSWORD string
	JAEGER_ENDPOINT   string
	MEMCACHED_HOST    string
}

func LoadConfig(configPath string) (*EnvConfig, error) {
	var (
		config EnvConfig
		err    error
	)

	viper.SetConfigFile(configPath)
	if err = viper.ReadInConfig(); err == nil {
		config = EnvConfig{
			HTTP_ADDRESS:  viper.GetString("HTTP_ADDRESS"),
			MAX_LIMITER:   viper.GetFloat64("MAX_LIMITER"),
			DB_DRIVER:     viper.GetString("DB_DRIVER"),
			DATABASE_NAME: viper.GetString("DB_DATASOURCE"),
		}
		return &config, nil
	}
	return nil, err
}
