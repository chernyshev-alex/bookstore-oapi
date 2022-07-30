package env

import (
	"fmt"

	"github.com/spf13/viper"
)

type DbConfig struct {
	Driver string
	Host   string
	Port   int
	User   string `mapstructure:"username"`
	Pass   string `mapstructure:"password"`
	Dbname string
}
type ServerConfig struct {
	Host           string `mapstructure:"hostname"`
	Port           int    `mapstructure:"port"`
	JaegerEndPoint string `mapstructure:"jaeger"`
	MemcachedHost  string
	Prometeus      string `mapstructure:"prometeus"`
	HttpLimiter    float64
}
type Config struct {
	Db     DbConfig     `mapstructure:"database"`
	Server ServerConfig `mapstructure:"server"`
}

func LoadConfig(configPath string) (c *Config, e error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := v.Unmarshal(&c); err != nil {
		fmt.Printf("couldn't read config: %s", err)
	}
	return c, nil
}
