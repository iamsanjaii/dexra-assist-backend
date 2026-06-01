package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	pdf "github.com/dslipak/pdf"

	chroma "github.com/amikos-tech/chroma-go/pkg/api/v2"
	"github.com/amikos-tech/chroma-go/pkg/embeddings"
	"github.com/dexra/backend/internal/config"
	"github.com/dexra/backend/internal/database"
	"github.com/dexra/backend/internal/models"
	"github.com/dexra/backend/internal/repositories"
	"github.com/dexra/backend/internal/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// geminiEmbedRequest / Response for the v1 REST endpoint
type geminiEmbedRequest struct {
	Model   string           `json:"model"`
	Content geminiEmbedPart  `json:"content"`
}
type geminiEmbedPart struct {
	Parts []geminiTextPart `json:"parts"`
}
type geminiTextPart struct {
	Text string `json:"text"`
}
type geminiEmbedResponse struct {
	Embedding struct {
		Values []float32 `json:"values"`
	} `json:"embedding"`
}

// embedText calls the Gemini v1 REST API using gemini-embedding-001 (3072-dim).
// The genai Go SDK uses v1beta by default which has a different model availability.
func embedText(ctx context.Context, apiKey, text string) ([]float32, error) {
	url := "https://generativelanguage.googleapis.com/v1/models/gemini-embedding-001:embedContent?key=" + apiKey

	reqBody := geminiEmbedRequest{
		Model: "models/gemini-embedding-001",
		Content: geminiEmbedPart{
			Parts: []geminiTextPart{{Text: text}},
		},
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gemini embed API returned %d: %s", resp.StatusCode, string(raw))
	}

	var result geminiEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Embedding.Values, nil
}

// extractText extracts plain text from a file based on its extension.
func extractText(path, ext string) (string, error) {
	switch strings.ToLower(ext) {
	case ".pdf":
		return extractPDFText(path)
	default:
		// For .txt, .csv and other text files: read and sanitize UTF-8
		raw, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return sanitizeUTF8(string(raw)), nil
	}
}

// extractPDFText uses dslipak/pdf to pull plain text from a PDF file.
func extractPDFText(path string) (string, error) {
	r, err := pdf.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}
	var buf strings.Builder
	totalPages := r.NumPage()
	for i := 1; i <= totalPages; i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}
		content, err := p.GetPlainText(nil)
		if err != nil {
			// skip pages that can't be read
			continue
		}
		buf.WriteString(content)
		buf.WriteByte('\n')
	}
	return sanitizeUTF8(buf.String()), nil
}

// sanitizeUTF8 replaces invalid UTF-8 byte sequences with the replacement char.
func sanitizeUTF8(s string) string {
	if utf8.ValidString(s) {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			b.WriteRune('\uFFFD')
		} else {
			b.WriteRune(r)
		}
		i += size
	}
	return b.String()
}

// chunkText splits text into overlapping chunks of ~1000 runes.
func chunkText(text string, chunkSize int) []string {
	runes := []rune(text)
	var chunks []string
	for i := 0; i < len(runes); i += chunkSize {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunk := strings.TrimSpace(string(runes[i:end]))
		if chunk != "" {
			chunks = append(chunks, chunk)
		}
	}
	return chunks
}

func UploadDocument(file *multipart.FileHeader) (*models.Document, error) {
	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		return nil, err
	}

	ext := filepath.Ext(file.Filename)
	uniqueFilename := uuid.New().String() + ext
	storagePath := filepath.Join("uploads", uniqueFilename)

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(storagePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, err
	}

	doc := &models.Document{
		Filename:         file.Filename,
		FileType:         ext,
		StoragePath:      storagePath,
		ProcessingStatus: "processing",
		ChunkCount:       0,
		CreatedAt:        time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := repositories.CreateDocument(ctx, doc); err != nil {
		return nil, err
	}

	// Trigger async processing
	go processDocument(doc.ID.Hex(), storagePath, ext)

	return doc, nil
}

func processDocument(id, path, ext string) {
	utils.Logger.Info("Started processing document", zap.String("id", id))

	// 1. Extract text from the file (PDF-aware, UTF-8 safe)
	content, err := extractText(path, ext)
	if err != nil {
		utils.Logger.Error("Failed to extract text from document", zap.Error(err))
		repositories.UpdateDocumentStatus(context.Background(), id, "failed", 0)
		return
	}

	content = strings.TrimSpace(content)
	if content == "" {
		utils.Logger.Error("Extracted content is empty — cannot process document")
		repositories.UpdateDocumentStatus(context.Background(), id, "failed", 0)
		return
	}

	// 2. Chunk the extracted text (rune-safe, 1000 runes per chunk)
	chunks := chunkText(content, 1000)
	utils.Logger.Info("Chunked document", zap.String("id", id), zap.Int("chunks", len(chunks)))

	// 3. Generate Embeddings & Store in Chroma
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Write total chunk count upfront so the frontend knows the denominator (processed_chunks stays at 0)
	if err := repositories.SetTotalChunkCount(context.Background(), id, len(chunks)); err != nil {
		utils.Logger.Error("Failed to write total chunk count", zap.Error(err))
	}

	collection, err := database.GetKnowledgeCollection(ctx)
	if err != nil {
		utils.Logger.Error("Failed to get Chroma collection", zap.Error(err))
		repositories.UpdateDocumentStatus(context.Background(), id, "failed", 0)
		return
	}

	var ids []chroma.DocumentID
	var texts []string
	var embs []embeddings.Embedding
	var metadatas []chroma.DocumentMetadata

	for i, chunk := range chunks {
		// Use v1 REST API — the genai Go SDK uses v1beta which doesn't support text-embedding-004
		values, err := embedText(ctx, config.AppConfig.GeminiAPIKey, chunk)
		if err != nil {
			utils.Logger.Error("Failed to generate embedding", zap.Int("chunk", i), zap.Error(err))
			continue
		}

		ids = append(ids, chroma.DocumentID(id+"-"+uuid.New().String()[:8]))
		texts = append(texts, chunk)
		embs = append(embs, embeddings.NewEmbeddingFromFloat32(values))

		metaMap := map[string]interface{}{
			"document_id": id,
			"chunk_index": i,
		}
		meta, _ := chroma.NewDocumentMetadataFromMap(metaMap)
		metadatas = append(metadatas, meta)

		// Increment live progress counter after each successful embedding
		_ = repositories.IncrementProcessedChunks(context.Background(), id)
	}

	if len(ids) > 0 {
		err = collection.Add(ctx,
			chroma.WithIDs(ids...),
			chroma.WithTexts(texts...),
			chroma.WithEmbeddings(embs...),
			chroma.WithMetadatas(metadatas...),
		)
		if err != nil {
			utils.Logger.Error("Failed to insert into ChromaDB", zap.Error(err))
			repositories.UpdateDocumentStatus(context.Background(), id, "failed", 0)
			return
		}
	}

	// 4. Update status to ready with final chunk count
	err = repositories.UpdateDocumentStatus(context.Background(), id, "ready", len(chunks))
	if err != nil {
		utils.Logger.Error("Failed to update document status", zap.Error(err))
	} else {
		utils.Logger.Info("Finished processing document", zap.String("id", id), zap.Int("chunks_stored", len(ids)))
	}
}

func GetDocuments(page, limit int) ([]models.Document, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	skip := (page - 1) * limit
	return repositories.GetDocuments(ctx, nil, int64(limit), int64(skip))
}

func DeleteDocument(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Optionally delete file from disk and ChromaDB embeddings
	return repositories.DeleteDocument(ctx, id)
}
