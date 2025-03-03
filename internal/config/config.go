package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	Auth     Auth     `yaml:"auth"`
}

type Server struct {
	Host string `yaml:"host" env:"HOST" env-required:"true"`
	Port string `yaml:"port" env:"PORT" env-required:"true"`
}

type Database struct {
	Path string `yaml:"path" env:"DB_PATH" env-required:"true"`
}

type Auth struct {
	Password string `yaml:"password" env:"AUTH_PASSWORD" env-required:"true"`
	Secret   string `yaml:"secret" env:"AUTH_SECRET" env-required:"true"`
}

// Загружаем конфиг из файла и переопределяем переменными окружения
func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	var cfg Config

	// Читаем конфигурацию из файла
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	log.Printf("HOST: %s", cfg.Server.Host)
	log.Printf("PORT: %s", cfg.Server.Port)
	log.Printf("DB_PATH: %s", cfg.Database.Path)
	log.Printf("PASSWORD: %s", cfg.Auth.Password)
	log.Printf("SECRET: %s", cfg.Auth.Secret)

	return &cfg
}
