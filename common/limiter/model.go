package limiter

type RedisConfig struct {
	Enable            bool         `mapstructure:"Enable"`
	RedisAddr         string       `mapstructure:"RedisAddr"` // host:port
	RedisPassword     string       `mapstructure:"RedisPassword"`
	RedisDB           int          `mapstructure:"RedisDB"`
	Timeout           int          `mapstructure:"Timeout"`
	Expiry            int          `mapstructure:"Expiry"` // second
	ReportLimitConfig *ReportLimit `mapstructure:"ReportLimitConfig"`
}

type ReportLimit struct {
	EnableReporting bool   `mapstructure:"EnableReporting"`
	BotToken        string `mapstructure:"BotToken"`
	ChatID          string `mapstructure:"ChatID"`
}
