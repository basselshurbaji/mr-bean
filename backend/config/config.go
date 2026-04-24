package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port   string
	DB     DBConfig
	Auth   AuthConfig
	Mailer MailerConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type AuthConfig struct {
	JWTSecret     string
	JWTExpiry     time.Duration
	RefreshExpiry time.Duration
}

type MailerConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func (d DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

func Load() Config {
	secret := getEnv("JWT_SECRET", "mr_bean")

	return Config{
		Port: getEnv("PORT", "8080"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "mr_bean"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Auth: AuthConfig{
			JWTSecret:     secret,
			JWTExpiry:     getMinutes("JWT_EXPIRY", 1),
			RefreshExpiry: getMinutes("REFRESH_EXPIRY", 1440),
		},
		Mailer: MailerConfig{
			Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:     getEnv("SMTP_PORT", "587"),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", ""),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// getMinutes reads an integer number of minutes from an env var.
// fallbackMinutes is used when the var is unset or unparseable.
func getMinutes(key string, fallbackMinutes int) time.Duration {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return time.Duration(n) * time.Minute
		}
	}
	return time.Duration(fallbackMinutes) * time.Minute
}
