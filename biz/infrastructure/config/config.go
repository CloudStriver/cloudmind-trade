package config

import (
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"os"
)

type Config struct {
	service.ServiceConf
	ListenOn string
	Mongo    struct {
		URL string
		DB  string
	}
	CacheConf cache.CacheConf
	Redis     *redis.RedisConf
	Balance   struct {
		DefaultFlow   int64
		DefaultMemory int64
		DefaultPoint  int64
	}
}

func NewConfig() (*Config, error) {
	c := new(Config)
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "etc/config.yaml"
	}
	err := conf.Load(path, c)
	if err != nil {
		return nil, err
	}
	err = c.SetUp()
	if err != nil {
		return nil, err
	}
	return c, nil
}
