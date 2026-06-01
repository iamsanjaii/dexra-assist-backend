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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"stats": gin.H{
			"totalDocuments":         totalDocs,
			"totalQAPairs":           totalQA,
			"totalConversations":     totalConv,
			"activeKnowledgeSources": activeSources,
		},
	})
}
