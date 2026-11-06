package config

type Config struct {
	Host          string   `env:"HOST"`
	Port          int      `env:"PORT" envDefault:"3000"`
	IsDevelopment bool     `env:"IS_DEVELOPMENT"`
	ProxyHeader   string   `env:"PROXY_HEADER"`
	Database      Database `envPrefix:"DB_"`
	LogFields     []string `env:"LOG_FIELDS" envSeparator:","`
}

type Database struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Name     string `env:"NAME"`
}
