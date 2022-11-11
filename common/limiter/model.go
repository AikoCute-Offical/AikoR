package limiter

type RedisConfig struct {
	RedisLimit    int    `mapstructure:"Limit"`
	RedisAddr     string `mapstructure:"RedisAddr"` // host:port
	RedisPassword string `mapstructure:"RedisPassword"`
	RedisDB       int    `mapstructure:"RedisDB"`
	Timeout       int    `mapstructure:"Timeout"`
	Expiry        int    `mapstructure:"Expiry"` // minute
}
