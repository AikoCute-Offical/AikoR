package limiter

type GlobalDeviceLimitConfig struct {
	RedisEnable   bool   `mapstructure:"RedisEnable"`
	RedisAddr     string `mapstructure:"RedisAddr"` // host:port
	RedisPassword string `mapstructure:"RedisPassword"`
	RedisDB       int    `mapstructure:"RedisDB"`
	RedisTimeout  int    `mapstructure:"RedisTimeout"`
	Expiry        int    `mapstructure:"Expiry"` // minute
}
