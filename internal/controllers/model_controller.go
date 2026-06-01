package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAvailableModels(c *gin.Context) {
	models := map[string][]string{
		"google":     {"gemini-1.5-flash", "gemini-1.5-pro", "gemini-2.0-flash"},
		"openrouter": {
			"anthropic/claude-3.5-sonnet",
			"meta-llama/llama-3.1-8b-instruct",
			"meta-llama/llama-3.1-70b-instruct",
			"google/gemini-pro-1.5",
		},
	}
	c.JSON(http.StatusOK, models)
}
