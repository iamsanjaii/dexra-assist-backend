package repositories

import (
	"context"

	"github.com/dexra/backend/internal/database"
	"github.com/dexra/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateChatSession(ctx context.Context, session *models.ChatSession) error {
	collection := database.GetCollection("chat_sessions")
	res, err := collection.InsertOne(ctx, session)
	if err == nil {
		session.ID = res.InsertedID.(primitive.ObjectID)
	}
	return err
}

func GetChatSessions(ctx context.Context, userID string) ([]models.ChatSession, error) {
	collection := database.GetCollection("chat_sessions")
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	opts := options.Find().SetSort(bson.M{"created_at": -1})
	cursor, err := collection.Find(ctx, bson.M{"user_id": uid}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []models.ChatSession
	if err = cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

func GetChatSessionByID(ctx context.Context, id string) (*models.ChatSession, error) {
	collection := database.GetCollection("chat_sessions")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var session models.ChatSession
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&session)
	return &session, err
}

func CreateChatMessage(ctx context.Context, msg *models.ChatMessage) error {
	collection := database.GetCollection("chat_messages")
	res, err := collection.InsertOne(ctx, msg)
	if err == nil {
		msg.ID = res.InsertedID.(primitive.ObjectID)
	}
	return err
}

func GetChatMessages(ctx context.Context, sessionID string) ([]models.ChatMessage, error) {
	collection := database.GetCollection("chat_messages")
	objID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		return nil, err
	}

	opts := options.Find().SetSort(bson.M{"created_at": 1})
	cursor, err := collection.Find(ctx, bson.M{"session_id": objID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var msgs []models.ChatMessage
	if err = cursor.All(ctx, &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}
