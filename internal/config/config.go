package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port               string `mapstructure:"PORT"`
	MongoURI           string `mapstructure:"MONGO_URI"`
	MongoDbName        string `mapstructure:"MONGO_DB_NAME"`
	ChromaURL          string `mapstructure:"CHROMA_URL"`
	ChromaAPIKey       string `mapstructure:"CHROMA_API_KEY"`
	ChromaTenantID     string `mapstructure:"CHROMA_TENANT_ID"`
	ChromaDatabase     string `mapstructure:"CHROMA_DATABASE"`
	JWTSecret          string `mapstructure:"JWT_SECRET"`
	GeminiAPIKey       string `mapstructure:"GEMINI_API_KEY"`
	OpenRouterAPIKey   string `mapstructure:"OPENROUTER_API_KEY"`
	ModelProvider      string `mapstructure:"MODEL_PROVIDER"`
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL  string `mapstructure:"GOOGLE_REDIRECT_URL"`
	FrontendURL        string `mapstructure:"FRONTEND_URL"`
	ChatbotURL         string `mapstructure:"CHATBOT_URL"`
	CookieDomain       string `mapstructure:"COOKIE_DOMAIN"`
	SecureCookies      bool   `mapstructure:"SECURE_COOKIES"`
}

var AppConfig *Config

// LoadConfig loads environment variables from the .env file or environment
func LoadConfig() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../..")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found or couldn't be loaded: %v", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	// Defaults
	if AppConfig.Port == "" {
		AppConfig.Port = "8080"
	}
	if AppConfig.ChatbotURL == "" {
		AppConfig.ChatbotURL = "http://localhost:3001"
	}
}
