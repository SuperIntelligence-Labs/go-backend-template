package config

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// Load reads configuration from .env file and environment variables.
func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found, falling back to environment variables")
	}

	viper.AutomaticEnv()

	cfg := &Config{
		Server: ServerConfig{
			Host: viper.GetString("SERVER_HOST"),
			Port: viper.GetString("SERVER_PORT"),
			Env:  viper.GetString("SERVER_ENV"),
		},
		Log: LogConfig{
			Level: viper.GetString("LOG_LEVEL"),
		},
		JWT: JWTConfig{
			ATSecret:    viper.GetString("JWT_AT_SECRET"),
			ATExpiresIn: viper.GetInt("JWT_AT_EXPIRES_IN"),
			RTSecret:    viper.GetString("JWT_RT_SECRET"),
			RTExpiresIn: viper.GetInt("JWT_RT_EXPIRES_IN"),
		},
		Db: DatabaseConfig{
			Host:            viper.GetString("DB_HOST"),
			Port:            viper.GetString("DB_PORT"),
			User:            viper.GetString("DB_USER"),
			Password:        viper.GetString("DB_PASSWORD"),
			Name:            viper.GetString("DB_NAME"),
			SSLMode:         viper.GetString("DB_SSL_MODE"),
			MaxOpenConns:    getIntWithDefault("DB_MAX_OPEN_CONNS", 10),
			MaxIdleConns:    getIntWithDefault("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getIntWithDefault("DB_CONN_MAX_LIFETIME", 60),
		},
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, ParseValidationErrors(err)
	}

	return cfg, nil
}

// getIntWithDefault returns the int value for the key or the default if not set
func getIntWithDefault(key string, defaultValue int) int {
	if viper.IsSet(key) {
		return viper.GetInt(key)
	}
	return defaultValue
}
