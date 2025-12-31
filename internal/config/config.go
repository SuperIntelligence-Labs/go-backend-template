package config

type Config struct {
	Server ServerConfig   `validate:"required"`
	Log    LogConfig      `validate:"required"`
	JWT    JWTConfig      `validate:"required"`
	Db     DatabaseConfig `validate:"required"`
}

type ServerConfig struct {
	Host string `validate:"required"`
	Port string `validate:"required,numeric"`
	Env  string `validate:"required,oneof=development production"`
}

type LogConfig struct {
	Level string `validate:"required,oneof=debug info warn error"`
}

type JWTConfig struct {
	ATSecret    string `validate:"required"`
	ATExpiresIn int    `validate:"min=1"`
	RTSecret    string `validate:"required"`
	RTExpiresIn int    `validate:"min=1"`
}

type DatabaseConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	Name     string `validate:"required"`
	SSLMode  string `validate:"required"`
}

