package services

import (
	"context"

	"github.com/dexra/backend/internal/models"
	"github.com/dexra/backend/internal/repositories"
)

func GetAIConfig() (*models.AIConfig, error) {
	return repositories.GetAIConfig(context.Background())
}

func UpdateAIConfig(config *models.AIConfig) error {
	return repositories.UpdateAIConfig(context.Background(), config)
}
