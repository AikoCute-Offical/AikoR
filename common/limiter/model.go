package limiter

type RedisConfig struct {
	RedisEnable   bool   `mapstructure:"RedisEnable"`   // if true, use redis to store online ip
	RedisAddr     string `mapstructure:"RedisAddr"`     // host:port
	RedisPassword string `mapstructure:"RedisPassword"` // redis password
	RedisDB       int    `mapstructure:"RedisDB"`       // redis db
	RedisTimeout  int    `mapstructure:"RedisTimeout"`  // second
	Expiry        int    `mapstructure:"Expiry"`        // second
}
