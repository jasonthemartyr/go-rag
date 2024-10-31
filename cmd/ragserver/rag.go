package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/milosgajdos/go-embeddings/ollama"
	"github.com/qdrant/go-client/qdrant"
)

const (
	vdbHost            = "localhost"
	vdbPort            = 6334
	collectionName     = "documents"
	dimensionEmbedding = 4096
	defaultOllamaURL   = "http://localhost:11434/api/chat"
)

type RAGHandler struct {
	QdrantClient *qdrant.Client
	OllamaClient *ollama.Client
	// OllamaClient   *api.Client //TODO: swap for official llama SDK
	CollectionName string
}

func NewRAGClients() (*RAGHandler, error) {
	qdClient, err := qdrant.NewClient(&qdrant.Config{
		Host: vdbHost,
		Port: vdbPort,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to qdrant: %w", err)
	}
	ollamaClient := ollama.NewClient()

	// ollamaClient, err := api.ClientFromEnvironment() //TODO: swap for official llama SDK
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to connect to qdrant: %w", err)
	// }
	return &RAGHandler{
		QdrantClient:   qdClient,
		OllamaClient:   ollamaClient,
		CollectionName: collectionName, //TODO: add conditional for default collection name
	}, nil
}

func (r *RAGHandler) CreateCollection() error {
	ctx := context.Background()
	collectionExists, err := r.QdrantClient.CollectionExists(ctx, r.CollectionName)
	if err != nil {
		return err
	}
	if !collectionExists {
		err := r.QdrantClient.CreateCollection(ctx, &qdrant.CreateCollection{
			CollectionName: r.CollectionName,
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     dimensionEmbedding,
				Distance: qdrant.Distance_Cosine,
			}),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func splitIntoChunks(text string, chunkSize int) []string {
	var chunks []string
	words := strings.Fields(text)
	var chunk []string
	length := 0

	for _, word := range words {
		if length+len(word) > chunkSize {
			chunks = append(chunks, strings.Join(chunk, " "))
			chunk = []string{word}
			length = len(word)
		} else {
			chunk = append(chunk, word)
			length += len(word) + 1
		}
	}

	if len(chunk) > 0 {
		chunks = append(chunks, strings.Join(chunk, " "))
	}

	return chunks
}

func (r *RAGHandler) ProcessDocuments(docs []Document, batchSize int) ([]*qdrant.PointStruct, error) {
	var points []*qdrant.PointStruct

	for i := 0; i < len(docs); i += batchSize {
		end := min(i+batchSize, len(docs))
		batch := docs[i:end]

		for j, doc := range batch {
			chunks := splitIntoChunks(doc.Content, 1000)

			for k, chunk := range chunks {

				// START ###############################
				//TODO: swap for official llama SDK
				// embReq := &api.EmbedRequest{
				// 	Model: "llama3",
				// 	Input: chunk,
				// }
				// embs, err := r.OllamaClient.Embed(context.Background(), embReq)
				// if err != nil {
				// 	return nil, fmt.Errorf("embedding error: %w", err)
				// }
				// embedVals := embs.Embeddings[0]

				embReq := &ollama.EmbeddingRequest{
					Prompt: chunk,
					Model:  "llama3",
				}

				embs, err := r.OllamaClient.Embed(context.Background(), embReq)
				if err != nil {
					return nil, fmt.Errorf("embedding error: %w", err)
				}
				embedVals := make([]float32, len(embs[0].Vector)) // Use the Vector field
				for i, val := range embs[0].Vector {
					embedVals[i] = float32(val)
				}
				// END ###############################
				point := &qdrant.PointStruct{
					Id:      qdrant.NewIDNum(uint64(i*batchSize + j*len(chunks) + k)),
					Vectors: qdrant.NewVectors(embedVals...),
					Payload: qdrant.NewValueMap(map[string]any{
						"text":     chunk,
						"metadata": doc.Metadata,
					}),
				}
				points = append(points, point)
			}
		}
	}
	return points, nil
}

func (r *RAGHandler) AddDocuments(points []*qdrant.PointStruct) (*qdrant.UpdateResult, error) {
	operationInfo, err := r.QdrantClient.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: r.CollectionName,
		Points:         points,
	},
	)
	if err != nil {
		return nil, fmt.Errorf("AddDocuments Upsert error: %w", err)
	}
	return operationInfo, nil

}

// func (r *RAGHandler) SearchDocuments(searchText string) ([]*qdrant.ScoredPoint, error) {
func (r *RAGHandler) SearchDocuments(searchText string) (string, error) {
	// START ###############################
	//TODO: swap for official llama SDK
	// embReq := &api.EmbedRequest{
	// 	Model: "llama3",
	// 	Input: searchText,
	// }
	// embs, err := r.OllamaClient.Embed(context.Background(), embReq)
	// if err != nil {
	// 	return "", fmt.Errorf("embedding error: %w", err)
	// }
	// searchVals := embs.Embeddings[0]

	embReq := &ollama.EmbeddingRequest{
		Prompt: searchText,
		Model:  "llama3",
	}
	embs, err := r.OllamaClient.Embed(context.Background(), embReq)
	if err != nil {
		return "", fmt.Errorf("search embedding error: %w", err)
	}

	searchVals := make([]float32, len(embs[0].Vector))
	for i, val := range embs[0].Vector {
		searchVals[i] = float32(val)
	}
	// END ###############################

	searchResults, err := r.QdrantClient.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: r.CollectionName,
		Query:          qdrant.NewQuery(searchVals...),
		WithPayload:    qdrant.NewWithPayload(true), //explicitly return payload otherwise its empty
		// WithPayload: &qdrant.WithPayloadSelector{},
		// Filter: &qdrant.Filter{ //TODO: figure out best filtering method
		// 	Must: []*qdrant.Condition{
		// 		qdrant.NewMatch("city", "London"),
		// 	},
		// },
	})
	if err != nil {
		return "", fmt.Errorf("searching error: %w", err)
	}

	// fmt.Printf("Found %d matches\n", len(searchResults))
	var context string
	for _, point := range searchResults {
		// fmt.Printf("Score: %.4f\n", point.Score)
		// fmt.Printf("ID: %v\n", point.Id)
		// fmt.Printf("Payload: %v\n\n", point.Payload)
		if payload, ok := point.Payload["text"]; ok {
			text := payload.GetStringValue() // or GetKind() to see what type it is
			// fmt.Printf("Text: %s\n", text)
			context += text + "\n"
		}

		// if metaValue, ok := point.Payload["metadata"]; ok {
		// 	meta := metaValue.GetStructValue()
		// 	fmt.Printf("Metadata: %+v\n", meta)
		// }

		// if point.Vectors != nil {
		// 	x := point.GetVectors()
		// 	fmt.Println("Vector length: ", x)
		// }
		// fmt.Printf("Version: %d\n", point.Version)
		// fmt.Println("---")

	}
	prompt := fmt.Sprintf(`Use the following context to answer the question. 
	Context: %s
	Question: %s
	Answer based on the context provided:`, context, searchText)

	msg := ChatMessage{
		Role:    "user",
		Content: prompt,
	}

	req := ChatRequest{
		Model:    "llama3",
		Stream:   false,
		Messages: []ChatMessage{msg},
	}
	response, err := talkToOllama(defaultOllamaURL, req)
	if err != nil {
		return "", fmt.Errorf("generate response: %w", err)
	}

	return response.Message.Content, nil
}

func talkToOllama(url string, ollamaReq ChatRequest) (*ChatResponse, error) {
	js, err := json.Marshal(&ollamaReq)
	if err != nil {
		return nil, err
	}
	// fmt.Println(ollamaReq)
	client := http.Client{}
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(js))
	if err != nil {
		return nil, err
	}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	ollamaResp := ChatResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&ollamaResp)
	return &ollamaResp, err
}
