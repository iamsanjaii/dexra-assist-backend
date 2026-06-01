package repositories

import (
	"context"

	"github.com/dexra/backend/internal/database"
	"github.com/dexra/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateUser inserts a new user
func CreateUser(ctx context.Context, user *models.User) error {
	collection := database.GetCollection("users")
	res, err := collection.InsertOne(ctx, user)
	if err == nil {
		user.ID = res.InsertedID.(primitive.ObjectID)
	}
	return err
}

// GetUserByEmail finds a user by email
func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	collection := database.GetCollection("users")
	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return &user, err
}

// GetUserByID finds a user by their MongoDB ObjectID string
func GetUserByID(ctx context.Context, id string) (*models.User, error) {
	collection := database.GetCollection("users")
	var user models.User
	
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	return &user, err
}

// UpdateUser updates an existing user
func UpdateUser(ctx context.Context, user *models.User) error {
	collection := database.GetCollection("users")
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{
			"name":    user.Name,
			"picture": user.Picture,
		}},
	)
	return err
}
