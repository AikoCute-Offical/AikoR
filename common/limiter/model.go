package limiter

import "github.com/go-redis/redis/v8"

type RedisConfig struct {
	RedisEnable   bool   `mapstructure:"RedisEnable"`
	RedisAddr     string `mapstructure:"RedisAddr"` // host:port
	RedisPassword string `mapstructure:"RedisPassword"`
	RedisDB       int    `mapstructure:"RedisDB"`
	RedisTimeout  int    `mapstructure:"RedisTimeout"`
	Expiry        int    `mapstructure:"Expiry"` // minute
	R             *redis.Client
}
