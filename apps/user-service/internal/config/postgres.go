package config

type Postgres struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	User	 string `koanf:"user"`
	Password string `koanf:"password"`
	Database string `koanf:"database"`
	SSLMode  string `koanf:"sslmode"`
}