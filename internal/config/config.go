package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	SMTPConfig  SMTPConfig
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func Load() *Config {
	//TODO: Read from yaml or secret server
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://loanuser:loanpassword@localhost:5432/loandb?sslmode=disable"),
		SMTPConfig: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.example.com"),
			Port:     587,
			Username: getEnv("SMTP_USERNAME", "noreply@example.com"),
			Password: getEnv("SMTP_PASSWORD", "smtp-password"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
