package database

import (
	"context"
	"time"

	"github.com/dexra/backend/internal/config"
	"github.com/dexra/backend/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

// ConnectMongoDB establishes a connection to MongoDB
func ConnectMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.AppConfig.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		utils.Logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		utils.Logger.Fatal("Failed to ping MongoDB", zap.Error(err))
	}

	MongoClient = client
	MongoDB = client.Database("dexra_assist")
	utils.Logger.Info("Connected to MongoDB successfully")
}

// GetCollection returns a MongoDB collection reference
func GetCollection(collectionName string) *mongo.Collection {
	return MongoDB.Collection(collectionName)
}
