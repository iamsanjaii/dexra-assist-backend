package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/dexra/backend/internal/models"
	"github.com/dexra/backend/internal/repositories"
	"github.com/dexra/backend/internal/services"
	"github.com/dexra/backend/internal/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		
		errorResponse := err.Error()
		if strings.Contains(errorResponse, "429") || strings.Contains(errorResponse, "RESOURCE_EXHAUSTED") || strings.Contains(errorResponse, "Resource exhausted") {
			errorResponse = "I'm receiving too many requests right now and hit a rate limit. Please wait a moment and try again."
		} else if strings.Contains(errorResponse, "security policies") || strings.Contains(errorResponse, "violates") {
			// Do not mask security violation errors
			errorResponse = err.Error()
		} else {
			errorResponse = "An internal error occurred while processing your request. Please try again later."
		}

		// Save the error response to the chat history so it persists
		if sessionObjID, parseErr := primitive.ObjectIDFromHex(req.SessionID); parseErr == nil {
			errorMsg := &models.ChatMessage{
				SessionID: sessionObjID,
				Role:      "assistant",
				Content:   errorResponse,
				CreatedAt: time.Now(),
			}
			repositories.CreateChatMessage(c.Request.Context(), errorMsg)
		}

		c.JSON(http.StatusOK, gin.H{
			"response":         errorResponse,
			"sources":          []string{},
			"confidence_score": 0.0,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response":         botMsg.Content,
		"sources":          botMsg.Sources,
		"confidence_score": 0.95, // Stub
	})
}
