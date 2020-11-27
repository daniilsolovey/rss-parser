package config

import "github.com/kovetskiy/ko"

type Database struct {
	Name      string `toml:"name" required:"true" env:"DATABASE_NAME"`
	Host      string `toml:"host" required:"true" env:"DATABASE_HOST"`
	Port      int    `toml:"port" required:"true" env:"DATABASE_PORT"`
	User      string `toml:"user" required:"true"`
	Password  string `toml:"password" required:"true"`
	TableName string `toml:"table_name" required:"true"`
}

type Config struct {
	Database Database `required:"true"`
	HTTPPort string   `toml:"http_port" required:"true"`
}

func Load(path string) (*Config, error) {
	config := &Config{}
	err := ko.Load(path, config, ko.RequireFile(false))
	if err != nil {
		return nil, err
	}

	return config, nil
}
