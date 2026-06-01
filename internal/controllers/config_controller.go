package controllers

import (
	"net/http"

	"github.com/dexra/backend/internal/models"
	"github.com/dexra/backend/internal/services"
	"github.com/gin-gonic/gin"
)

func GetAIConfig(c *gin.Context) {
	config, err := services.GetAIConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get AI config", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, config)
}

func UpdateAIConfig(c *gin.Context) {
	var config models.AIConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := services.UpdateAIConfig(&config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update AI config"})
		return
	}
	c.JSON(http.StatusOK, config)
}
