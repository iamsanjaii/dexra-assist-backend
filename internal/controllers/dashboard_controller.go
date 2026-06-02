package controllers

import (
	"net/http"

	"github.com/dexra/backend/internal/repositories"
	"github.com/gin-gonic/gin"
)

// GetDashboardStats returns real stats for the dashboard
func GetDashboardStats(c *gin.Context) {
	ctx := c.Request.Context()

	totalDocs, _ := repositories.GetTotalDocuments(ctx)
	totalQA, _ := repositories.GetTotalQAPairs(ctx)
	totalConv, _ := repositories.GetTotalConversations(ctx)
	activeSources, _ := repositories.GetActiveKnowledgeSources(ctx)
	
	aiStats, _ := repositories.GetAIAnalyticsStats(ctx)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"stats": gin.H{
			"totalDocuments":         totalDocs,
			"totalQAPairs":           totalQA,
			"totalConversations":     totalConv,
			"activeKnowledgeSources": activeSources,
		},
		"aiAnalytics": aiStats,
	})
}

// GetActivityFeed returns the recent activity logs
func GetActivityFeed(c *gin.Context) {
	ctx := c.Request.Context()

	logs, err := repositories.GetRecentActivityLogs(ctx)
	if err != nil {
		// Return empty array on error so UI doesn't crash
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"data":    []interface{}{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    logs,
	})
}
