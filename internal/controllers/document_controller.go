package controllers

import (
	"net/http"
	"strconv"

	"github.com/dexra/backend/internal/repositories"
	"github.com/dexra/backend/internal/services"
	"github.com/gin-gonic/gin"
)

func UploadDocument(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "No file uploaded"})
		return
	}

	doc, err := services.UploadDocument(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to upload document", "error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	var userIDStr string
	if userID != nil {
		userIDStr = userID.(string)
	}
	go repositories.LogActivity(userIDStr, "Uploaded Document", file.Filename)

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": gin.H{
		"document_id": doc.ID.Hex(),
		"status":      doc.ProcessingStatus,
	}})
}

func GetDocuments(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	docs, err := services.GetDocuments(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch documents"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": docs})
}

func DeleteDocument(c *gin.Context) {
	id := c.Param("id")
	
	if err := services.DeleteDocument(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Document deleted"})
}
