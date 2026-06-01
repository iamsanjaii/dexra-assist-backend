package repositories

import (
	"context"

	"github.com/dexra/backend/internal/database"
	"github.com/dexra/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAIConfig(ctx context.Context) (*models.AIConfig, error) {
	collection := database.GetCollection("configs")
	var config models.AIConfig
	err := collection.FindOne(ctx, bson.M{}).Decode(&config)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return default config if none exists
			return &models.AIConfig{
				Provider: "google",
				Model:    "gemini-1.5-flash",
			}, nil
		}
		return nil, err
	}
	return &config, nil
}

func UpdateAIConfig(ctx context.Context, config *models.AIConfig) error {
	collection := database.GetCollection("configs")
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, bson.M{}, bson.M{"$set": config}, opts)
	return err
}
