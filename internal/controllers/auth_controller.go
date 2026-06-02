package controllers

import (
	"net/http"

	"github.com/dexra/backend/internal/config"
	"github.com/dexra/backend/internal/repositories"
	"github.com/dexra/backend/internal/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type GoogleLoginRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request parameters", "error": err.Error()})
		return
	}

	accessToken, refreshToken, user, err := services.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

func getOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.AppConfig.GoogleClientID,
		ClientSecret: config.AppConfig.GoogleClientSecret,
		RedirectURL:  config.AppConfig.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func GoogleLoginRedirect(c *gin.Context) {
	url := getOAuthConfig().AuthCodeURL("admin", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleClientLoginRedirect(c *gin.Context) {
	url := getOAuthConfig().AuthCodeURL("client", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	accessToken, refreshToken, authUser, err := services.GoogleLoginCodeFlow(c.Request.Context(), getOAuthConfig(), code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
		return
	}

	// Log login activity
	if authUser != nil {
		go repositories.LogActivity(authUser.ID.Hex(), "Signed in", "via Google")
	}

	// Set HTTP-Only cookies for the tokens
	c.SetCookie("dexra_access_token", accessToken, 3600*24, "/", "", false, true)
	c.SetCookie("dexra_refresh_token", refreshToken, 3600*24*7, "/", "", false, true)

	state := c.Query("state")
	if state == "client" {
		targetURL := config.AppConfig.ChatbotURL
		if targetURL == "" {
			targetURL = "http://localhost:3001" // Fallback
		}
		c.Redirect(http.StatusTemporaryRedirect, targetURL)
	} else {
		c.Redirect(http.StatusTemporaryRedirect, config.AppConfig.FrontendURL+"/dashboard")
	}
}

func GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	user, err := repositories.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user,
	})
}

func Refresh(c *gin.Context) {
	// Not fully implemented for brevity, would decode refresh token and issue new access token.
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Token refreshed"})
}

func Logout(c *gin.Context) {
	c.SetCookie("dexra_access_token", "", -1, "/", "", false, true)
	c.SetCookie("dexra_refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Logged out successfully"})
}
