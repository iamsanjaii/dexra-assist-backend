package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/dexra/backend/internal/database"
	"github.com/dexra/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// LogActivity is a helper to quickly create activity logs
func LogActivity(userIDStr, action, item string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userID primitive.ObjectID
	if userIDStr != "" {
		if id, err := primitive.ObjectIDFromHex(userIDStr); err == nil {
			userID = id
		}
	}

	log := &models.ActivityLog{
		UserID:    userID,
		Action:    action,
		Item:      item,
		Time:      "Just now", // MVP static time or simple calculation later
		CreatedAt: time.Now(),
	}

	collection := database.GetCollection("activity_logs")
	collection.InsertOne(ctx, log)
}

// GetRecentActivityLogs retrieves the latest 10 logs
func GetRecentActivityLogs(ctx context.Context) ([]models.ActivityLog, error) {
	collection := database.GetCollection("activity_logs")

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	findOptions.SetLimit(10)

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []models.ActivityLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	// Update "Time" format on the fly for better UI if wanted, or just rely on what was saved
	for i, log := range logs {
		duration := time.Since(log.CreatedAt)
		if duration.Hours() > 24 {
			logs[i].Time = fmt.Sprintf("%d days ago", int(duration.Hours()/24))
		} else if duration.Hours() > 1 {
			logs[i].Time = fmt.Sprintf("%d hours ago", int(duration.Hours()))
		} else if duration.Minutes() > 1 {
			logs[i].Time = fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
		} else {
			logs[i].Time = "Just now"
		}
	}

	return logs, nil
}
