package repositories

import (
	"context"

	"github.com/dexra/backend/internal/database"
	"go.mongodb.org/mongo-driver/bson"
)

func GetTotalDocuments(ctx context.Context) (int64, error) {
	collection := database.GetCollection("documents")
	return collection.CountDocuments(ctx, bson.M{})
}

func GetTotalQAPairs(ctx context.Context) (int64, error) {
	collection := database.GetCollection("qa_pairs")
	return collection.CountDocuments(ctx, bson.M{})
}

func GetTotalConversations(ctx context.Context) (int64, error) {
	collection := database.GetCollection("chat_sessions")
	return collection.CountDocuments(ctx, bson.M{})
}

func GetActiveKnowledgeSources(ctx context.Context) (int64, error) {
	collection := database.GetCollection("documents")
	return collection.CountDocuments(ctx, bson.M{"processing_status": "ready"})
}
