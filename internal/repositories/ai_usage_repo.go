package repositories

import (
	"context"

	"github.com/dexra/backend/internal/database"
	"github.com/dexra/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateAIUsageLog(ctx context.Context, log *models.AIUsageLog) error {
	collection := database.GetCollection("ai_usage_logs")
	_, err := collection.InsertOne(ctx, log)
	return err
}

type DailyStat struct {
	Date     string `json:"name"`
	Tokens   int    `json:"tokens"`
	Requests int    `json:"requests"`
}

type AIAnalyticsStats struct {
	TotalTokens    int
	TotalRequests  int
	AvgLatency     int
	EstimatedCost  float64
	DailyStats     []DailyStat
}

func GetAIAnalyticsStats(ctx context.Context) (*AIAnalyticsStats, error) {
	collection := database.GetCollection("ai_usage_logs")

	// Get total requests
	totalRequests, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	// Calculate sums and averages
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":              nil,
				"totalTokens":      bson.M{"$sum": "$total_tokens"},
				"avgResponseTime":  bson.M{"$avg": "$response_time_ms"},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var result []struct {
		TotalTokens     int     `bson:"totalTokens"`
		AvgResponseTime float64 `bson:"avgResponseTime"`
	}

	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	cursor.Close(ctx)

	stats := &AIAnalyticsStats{
		TotalRequests: int(totalRequests),
		DailyStats:    make([]DailyStat, 0),
	}

	if len(result) > 0 {
		stats.TotalTokens = result[0].TotalTokens
		stats.AvgLatency = int(result[0].AvgResponseTime)
		stats.EstimatedCost = float64(stats.TotalTokens) / 1000000.0 * 1.0
	}

	// Fetch daily stats
	dailyPipeline := []bson.M{
		{
			"$group": bson.M{
				"_id": bson.M{
					"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$created_at"},
				},
				"tokens":   bson.M{"$sum": "$total_tokens"},
				"requests": bson.M{"$sum": 1},
			},
		},
		{
			"$sort": bson.M{"_id": 1},
		},
		{
			"$limit": 7, // last 7 days roughly
		},
	}

	dailyCursor, err := collection.Aggregate(ctx, dailyPipeline)
	if err == nil {
		var dailyResult []struct {
			ID       string `bson:"_id"`
			Tokens   int    `bson:"tokens"`
			Requests int    `bson:"requests"`
		}
		if err = dailyCursor.All(ctx, &dailyResult); err == nil {
			for _, r := range dailyResult {
				stats.DailyStats = append(stats.DailyStats, DailyStat{
					Date:     r.ID,
					Tokens:   r.Tokens,
					Requests: r.Requests,
				})
			}
		}
		dailyCursor.Close(ctx)
	}

	return stats, nil
}
