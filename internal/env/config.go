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
	Port           int
	JaegerEndPoint string `mapstructure:"jaeger"`
	MemcachedHost  string `mapstructure:"memcache"`
	Prometeus      string
	HttpLimiter    float64
}
type Config struct {
	Db     DbConfig     `mapstructure:"database"`
	Server ServerConfig `mapstructure:"server"`
}

func LoadConfig(configPath string) (c *Config, e error) {
	fmt.Printf("*** configPath *** %s \n", configPath)

	v := viper.New()
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("*** ReadInConfig *** %v\n", err)
		return nil, err
	}

	if err := v.Unmarshal(&c); err != nil {
		fmt.Printf("couldn't read config: %s\n", err)
	}
	return c, nil
}
