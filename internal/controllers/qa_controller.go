package controllers

import (
	"net/http"

	"github.com/dexra/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type QARequest struct {
	Question string `json:"question" binding:"required"`
	Answer   string `json:"answer" binding:"required"`
}

func CreateQAPair(c *gin.Context) {
	var req QARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Validation failed", "error": err.Error()})
		return
	}

	qa, err := services.CreateQAPair(req.Question, req.Answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create QA pair"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": qa})
}

func GetQAPairs(c *gin.Context) {
	qas, err := services.GetQAPairs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch QA pairs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": qas})
}

func UpdateQAPair(c *gin.Context) {
	id := c.Param("id")
	var req QARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Validation failed", "error": err.Error()})
		return
	}

	if err := services.UpdateQAPair(id, req.Question, req.Answer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update QA pair"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "QA pair updated"})
}

func DeleteQAPair(c *gin.Context) {
	id := c.Param("id")

	if err := services.DeleteQAPair(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete QA pair"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "QA pair deleted"})
}
