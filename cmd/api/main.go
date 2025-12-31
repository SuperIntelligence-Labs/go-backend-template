package main

import (
	"log"

	"github.com/SuperIntelligence-Labs/go-backend-template/internal/config"
	"github.com/SuperIntelligence-Labs/go-backend-template/internal/database"
	"github.com/SuperIntelligence-Labs/go-backend-template/internal/features/example"
	"github.com/SuperIntelligence-Labs/go-backend-template/internal/logger"
	"github.com/SuperIntelligence-Labs/go-backend-template/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config > %v", err)
	}

	config.SetEnv(cfg.Server.Env)
	logger.Init(cfg.Log.Level)

	// Database Setup
	db, err := database.NewDB(&cfg.Db)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Ensure database connection is closed on shutdown
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get sql.DB from GORM")
	}
	defer sqlDB.Close()

	logger.Info().Msg("Connected to database")

	// Auto migrate models
	err = db.AutoMigrate(&example.Item{})
	if err != nil {
		logger.Fatal().Err(err).Msg("Database migration failed")
	}
	logger.Info().Msg("Database migrated successfully")

	// Dependency Injection - Example Feature
	exampleRepo := example.NewRepository(db)
	exampleService := example.NewService(exampleRepo)
	exampleHandler := example.NewHandler(exampleService)

	// Server Setup
	srv := server.New()
	srv.RegisterRoutes(server.RoutesConfig{
		ExampleHandler: exampleHandler,
	})

	logger.Info().Str("port", cfg.Server.Port).Msg("Starting server")
	if err := srv.Start(":" + cfg.Server.Port); err != nil {
		logger.Fatal().Err(err).Msg("Server failed to start")
	}
}
