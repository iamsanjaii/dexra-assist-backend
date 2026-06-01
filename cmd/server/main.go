package main

import (
	"log"

	"github.com/dexra/backend/internal/config"
	"github.com/dexra/backend/internal/database"
	"github.com/dexra/backend/internal/routes"
	"github.com/dexra/backend/internal/services"
	"github.com/dexra/backend/internal/utils"
	"go.uber.org/zap"
)

func main() {
	// Initialize Config
	config.LoadConfig()

	// Initialize Logger
	utils.InitLogger()
	defer utils.SyncLogger()

	// Connect to Databases
	utils.Logger.Info("Connecting to MongoDB...")
	database.ConnectMongoDB()

	utils.Logger.Info("Connecting to ChromaDB...")
	database.ConnectChromaDB()

	// Seed default admin user
	utils.Logger.Info("Seeding Admin User...")
	if err := services.SeedAdminUser(); err != nil {
		utils.Logger.Error("Failed to seed admin user", zap.Error(err))
	} else {
		utils.Logger.Info("Admin User seeded successfully (admin@dexra.ai / password)")
	}

	// Setup Routes
	r := routes.SetupRouter()

	// Start Server
	port := config.AppConfig.Port
	utils.Logger.Info("Starting server on port " + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
