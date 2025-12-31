package config

// Config holds the complete application configuration.
type Config struct {
	Server ServerConfig   `validate:"required"`
	Log    LogConfig      `validate:"required"`
	JWT    JWTConfig      `validate:"required"`
	Db     DatabaseConfig `validate:"required"`
}

// ServerConfig defines HTTP server settings.
type ServerConfig struct {
	Host string `validate:"required"`
	Port string `validate:"required,numeric"`
	Env  string `validate:"required,oneof=development production"`
}

// LogConfig defines logging settings.
type LogConfig struct {
	Level string `validate:"required,oneof=debug info warn error"`
}

// JWTConfig defines JWT authentication settings.
type JWTConfig struct {
	ATSecret    string `validate:"required"`
	ATExpiresIn int    `validate:"min=1"`
	RTSecret    string `validate:"required"`
	RTExpiresIn int    `validate:"min=1"`
}

// DatabaseConfig defines PostgreSQL connection settings.
type DatabaseConfig struct {
	Host            string `validate:"required"`
	Port            string `validate:"required"`
	User            string `validate:"required"`
	Password        string `validate:"required"`
	Name            string `validate:"required"`
	SSLMode         string `validate:"required"`
	MaxOpenConns    int    `validate:"min=1"`
	MaxIdleConns    int    `validate:"min=1"`
	ConnMaxLifetime int    `validate:"min=1"` // in minutes
}

