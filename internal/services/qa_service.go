package services

import (
	"context"
	"time"

	"github.com/dexra/backend/internal/models"
	"github.com/dexra/backend/internal/repositories"
)

func CreateQAPair(question, answer string) (*models.QAPair, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	qa := &models.QAPair{
		Question:  question,
		Answer:    answer,
		CreatedAt: time.Now(),
	}

	err := repositories.CreateQAPair(ctx, qa)
	if err != nil {
		return nil, err
	}

	// TODO: Here we should also embed the QA pair and store in ChromaDB
	// StoreEmbeddingInChroma(...)

	return qa, nil
}

func GetQAPairs() ([]models.QAPair, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return repositories.GetQAPairs(ctx)
}

func UpdateQAPair(id, question, answer string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repositories.UpdateQAPair(ctx, id, question, answer)
	if err != nil {
		return err
	}

	// TODO: Update in ChromaDB as well
	return nil
}

func DeleteQAPair(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repositories.DeleteQAPair(ctx, id)
	if err != nil {
		return err
	}

	// TODO: Delete from ChromaDB as well
	return nil
}
