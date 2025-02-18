package config

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
	// In a real application, this would load from environment variables
	// For demonstration purposes, we're using hardcoded values
	return &Config{
		Port:        "9999",
		DatabaseURL: "postgres://username:password@localhost:5432/loandb?sslmode=disable",
		SMTPConfig: SMTPConfig{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "noreply@example.com",
			Password: "smtp-password",
		},
	}
}
