package routes

import (
	"github.com/dexra/backend/internal/controllers"
	"github.com/dexra/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	
	// Set a higher limit for multipart forms (50 MB) to allow large document uploads
	r.MaxMultipartMemory = 50 << 20 

	r.Use(gin.Logger())
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())

	// Serve the uploads directory statically
	r.Static("/uploads", "./uploads")

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", middleware.RateLimitMiddleware(), controllers.Login)
			auth.GET("/google/login", controllers.GoogleLoginRedirect)
			auth.GET("/client/google/login", controllers.GoogleClientLoginRedirect)
			auth.GET("/google/callback", controllers.GoogleCallback)
			auth.POST("/refresh", controllers.Refresh)
			auth.POST("/logout", controllers.Logout)

			// Protected auth route
			auth.GET("/me", middleware.AuthMiddleware(), controllers.GetMe)
		}

		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			dashboard := protected.Group("/dashboard")
			{
				dashboard.GET("/stats", controllers.GetDashboardStats)
				dashboard.GET("/activity", controllers.GetActivityFeed)
			}

			docs := protected.Group("/documents")
			{
				docs.POST("/upload", controllers.UploadDocument)
				docs.GET("", controllers.GetDocuments)
				docs.DELETE("/:id", controllers.DeleteDocument)
			}

			qa := protected.Group("/qa")
			{
				qa.POST("", controllers.CreateQAPair)
				qa.GET("", controllers.GetQAPairs)
				qa.PUT("/:id", controllers.UpdateQAPair)
				qa.DELETE("/:id", controllers.DeleteQAPair)
			}

			chat := protected.Group("/chat")
			{
				chat.POST("/session", controllers.CreateChatSession)
				chat.GET("/sessions", controllers.GetChatSessions)
				chat.GET("/history/:session_id", controllers.GetChatHistory)
				chat.POST("/query", middleware.RateLimitMiddleware(), controllers.HandleChatQuery)
			}

			config := protected.Group("/config")
			{
				config.GET("", controllers.GetAIConfig)
				config.PUT("", controllers.UpdateAIConfig)
			}

			models := protected.Group("/models")
			{
				models.GET("", controllers.GetAvailableModels)
			}
		}
	}

	return r
}
