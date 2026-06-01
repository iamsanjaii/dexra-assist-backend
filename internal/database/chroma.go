package database

import (
	"context"

	chroma "github.com/amikos-tech/chroma-go/pkg/api/v2"
	"github.com/dexra/backend/internal/config"
	"github.com/dexra/backend/internal/utils"
	"go.uber.org/zap"
)

var ChromaClient chroma.Client
var KnowledgeCollection *chroma.Collection

// ConnectChromaDB establishes a connection to ChromaDB Cloud
func ConnectChromaDB() {
	opts := []chroma.ClientOption{
		// Use the correct X-Chroma-Token header (not Authorization: Bearer)
		chroma.WithCloudAPIKey(config.AppConfig.ChromaAPIKey),
		// Set tenant and database from config
		chroma.WithDatabaseAndTenant(config.AppConfig.ChromaDatabase, config.AppConfig.ChromaTenantID),
	}

	client, err := chroma.NewCloudClient(opts...)
	if err != nil {
		utils.Logger.Fatal("Failed to connect to ChromaDB Cloud", zap.Error(err))
	}

	ChromaClient = client
	utils.Logger.Info("Connected to ChromaDB Cloud successfully")
}

// GetKnowledgeCollection returns or creates the knowledge_base collection
func GetKnowledgeCollection(ctx context.Context) (chroma.Collection, error) {
	if ChromaClient == nil {
		utils.Logger.Fatal("ChromaClient is not initialized")
	}

	// Use GetOrCreateCollection — atomically gets or creates, avoids nil-database race
	collection, err := ChromaClient.GetOrCreateCollection(ctx, "knowledge_base")
	if err != nil {
		return nil, err
	}

	return collection, nil
}

