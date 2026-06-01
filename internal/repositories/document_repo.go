package repositories

import (
	"context"

	"github.com/dexra/backend/internal/database"
	"github.com/dexra/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateDocument(ctx context.Context, doc *models.Document) error {
	collection := database.GetCollection("documents")
	res, err := collection.InsertOne(ctx, doc)
	if err == nil {
		doc.ID = res.InsertedID.(primitive.ObjectID)
	}
	return err
}

func GetDocuments(ctx context.Context, filter bson.M, limit int64, skip int64) ([]models.Document, error) {
	collection := database.GetCollection("documents")
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.M{"created_at": -1})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []models.Document
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

func GetDocumentByID(ctx context.Context, id string) (*models.Document, error) {
	collection := database.GetCollection("documents")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var doc models.Document
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&doc)
	return &doc, err
}

func UpdateDocumentStatus(ctx context.Context, id string, status string, chunks int) error {
	collection := database.GetCollection("documents")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"processing_status": status,
			"chunk_count":       chunks,
			"processed_chunks":  chunks, // sync processed to final count on completion
		},
	}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

// SetTotalChunkCount writes the total number of chunks discovered, without touching processed_chunks.
// Call this at the start of processing so the frontend knows the denominator.
func SetTotalChunkCount(ctx context.Context, id string, total int) error {
	collection := database.GetCollection("documents")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$set": bson.M{"chunk_count": total},
	})
	return err
}

// IncrementProcessedChunks atomically increments the processed_chunks counter
func IncrementProcessedChunks(ctx context.Context, id string) error {
	collection := database.GetCollection("documents")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$inc": bson.M{"processed_chunks": 1},
	})
	return err
}

func DeleteDocument(ctx context.Context, id string) error {
	collection := database.GetCollection("documents")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
