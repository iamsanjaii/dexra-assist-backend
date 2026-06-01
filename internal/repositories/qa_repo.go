package repositories

import (
	"context"

	"github.com/dexra/backend/internal/database"
	"github.com/dexra/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateQAPair(ctx context.Context, qa *models.QAPair) error {
	collection := database.GetCollection("qa_pairs")
	res, err := collection.InsertOne(ctx, qa)
	if err == nil {
		qa.ID = res.InsertedID.(primitive.ObjectID)
	}
	return err
}

func GetQAPairs(ctx context.Context) ([]models.QAPair, error) {
	collection := database.GetCollection("qa_pairs")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var qas []models.QAPair
	if err = cursor.All(ctx, &qas); err != nil {
		return nil, err
	}
	return qas, nil
}

func GetQAPairByID(ctx context.Context, id string) (*models.QAPair, error) {
	collection := database.GetCollection("qa_pairs")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var qa models.QAPair
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&qa)
	return &qa, err
}

func UpdateQAPair(ctx context.Context, id string, question string, answer string) error {
	collection := database.GetCollection("qa_pairs")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"question": question,
			"answer":   answer,
		},
	}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func DeleteQAPair(ctx context.Context, id string) error {
	collection := database.GetCollection("qa_pairs")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
