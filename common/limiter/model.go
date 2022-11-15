package limiter

type RedisConfig struct {
	RedisLimit    int    `mapstructure:"RedisLimit"`
	RedisAddr     string `mapstructure:"RedisAddr"` // host:port
	RedisPassword string `mapstructure:"RedisPassword"`
	RedisDB       int    `mapstructure:"RedisDB"`
	RedisTimeout  int    `mapstructure:"RedisTimeout"`
	Expiry        int    `mapstructure:"Expiry"` // minute
}
