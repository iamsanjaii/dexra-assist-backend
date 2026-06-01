package controllers

import (
	"net/http"

	"github.com/dexra/backend/internal/services"
	"github.com/dexra/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ChatSessionRequest struct {
	Title string `json:"title" binding:"required"`
}

type ChatQueryRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	Message   string `json:"message" binding:"required"`
}

func CreateChatSession(c *gin.Context) {
	var req ChatSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Validation failed", "error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	session, err := services.CreateChatSession(userID, req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create chat session"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": session})
}

func GetChatSessions(c *gin.Context) {
	userID := c.GetString("user_id")
	sessions, err := services.GetChatSessions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch chat sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": sessions})
}

func GetChatHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	sessionID := c.Param("session_id")
	messages, err := services.GetChatHistory(userID, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch chat history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": messages})
}

func HandleChatQuery(c *gin.Context) {
	var req ChatQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Validation failed", "error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	_, botMsg, err := services.HandleChatQuery(userID, req.SessionID, req.Message)
	if err != nil {
		utils.Logger.Error("HandleChatQuery failed", zap.String("session", req.SessionID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to process chat query", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response":         botMsg.Content,
		"sources":          botMsg.Sources,
		"confidence_score": 0.95, // Stub
	})
}
