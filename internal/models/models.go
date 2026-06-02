package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email" json:"email"`
	Name         string             `bson:"name" json:"name"`
	Picture      string             `bson:"picture" json:"picture"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}

type Document struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Filename         string             `bson:"filename" json:"filename"`
	FileType         string             `bson:"file_type" json:"file_type"`
	StoragePath      string             `bson:"storage_path" json:"storage_path"`
	ProcessingStatus string             `bson:"processing_status" json:"processing_status"` // e.g., "processing", "ready", "failed"
	ChunkCount       int                `bson:"chunk_count" json:"chunk_count"`
	ProcessedChunks  int                `bson:"processed_chunks" json:"processed_chunks"` // how many chunks have been embedded so far
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
}

type QAPair struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Question  string             `bson:"question" json:"question"`
	Answer    string             `bson:"answer" json:"answer"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type ChatSession struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Title     string             `bson:"title" json:"title"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type Source struct {
	DocumentID string `json:"document_id" bson:"document_id"`
	Name       string `json:"name" bson:"name"`
	Chunk      string `json:"chunk" bson:"chunk"`
}

type ChatMessage struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SessionID primitive.ObjectID `bson:"session_id" json:"session_id"`
	Role      string             `bson:"role" json:"role"` // "user" or "assistant"
	Content   string             `bson:"content" json:"content"`
	Sources   []Source           `bson:"sources,omitempty" json:"sources,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type AIConfig struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Provider string             `bson:"provider" json:"provider"`
	Model    string             `bson:"model" json:"model"`
}

type AIUsageLog struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SessionID        primitive.ObjectID `bson:"session_id" json:"session_id"`
	Query            string             `bson:"query" json:"query"`
	Model            string             `bson:"model" json:"model"`
	PromptTokens     int                `bson:"prompt_tokens" json:"prompt_tokens"`
	CompletionTokens int                `bson:"completion_tokens" json:"completion_tokens"`
	TotalTokens      int                `bson:"total_tokens" json:"total_tokens"`
	ResponseTimeMs   int64              `bson:"response_time_ms" json:"response_time_ms"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
}

type ActivityLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Action    string             `bson:"action" json:"action"`
	Item      string             `bson:"item" json:"item"`
	Time      string             `bson:"time" json:"time"` // formatted time (e.g., "Just now", "2 hours ago") to easily send to frontend MVP
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
