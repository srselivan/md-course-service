package config

import (
	"strings"

	"github.com/rs/zerolog"
)

const maskedSecret = "***"

// LogStartup writes all configuration values using .env variable names as field keys.
func (c *Config) LogStartup(log *zerolog.Logger) {
	kafkaBrokers := strings.Join(c.Kafka.Brokers, ",")

	log.Info().
		Str("LOG_LEVEL", c.Logger.Level).
		Str("LOG_FILE_PATH", c.Logger.FilePath).
		Str("HTTP_SERVER_ADDR", c.HTTPServer.Addr).
		Str("POSTGRES_HOST", c.Postgres.Host).
		Str("POSTGRES_PORT", c.Postgres.Port).
		Str("POSTGRES_USER", c.Postgres.User).
		Str("POSTGRES_PASSWORD", maskIfSet(c.Postgres.Password)).
		Str("POSTGRES_DB", c.Postgres.DBName).
		Str("POSTGRES_SSLMODE", c.Postgres.SSLMode).
		Str("POSTGRES_MIGRATIONS_PATH", c.Postgres.MigrationsPath).
		Str("KAFKA_BROKERS", kafkaBrokers).
		Msg("configuration loaded")
}

func maskIfSet(value string) string {
	if value == "" {
		return ""
	}
	return maskedSecret
}
