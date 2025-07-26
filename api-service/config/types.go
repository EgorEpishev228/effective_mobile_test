package config

type DNS struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"dbname"`
	Port     int    `mapstructure:"port"`
	SslMode  string `mapstructure:"sslmode"`
}

type Database struct {
	DNS DNS `mapstructure:"dns"`
}

type Server struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Config struct {
	Database Database `mapstructure:"database"`
	Server   Server   `mapstructure:"server"`
}

type FConfig struct {
	Prod Config `mapstructure:"prod"`
	Dev  Config `mapstructure:"dev,omitempty"`
	Test Config `mapstructure:"test,omitempty"`
}
