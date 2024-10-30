package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/milosgajdos/go-embeddings/ollama"
	"github.com/qdrant/go-client/qdrant"
)

const (
	vdbHost        = "localhost"
	vdbPort        = 6334
	collectionName = "documents"
	// dimensionEmbedding = 384
	dimensionEmbedding = 4096
)

type QdrantDB struct {
	Client         *qdrant.Client
	CollectionName string
}

//TODO: create ollama struct and return client

func NewQdrantDB() (*QdrantDB, error) {
	fmt.Println("connect to new DB")
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: vdbHost,
		Port: vdbPort,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to qdrant: %w", err)
	}
	return &QdrantDB{
		Client:         client,
		CollectionName: collectionName,
	}, nil
}

func (q *QdrantDB) CreateCollection() error {
	fmt.Println("create collection")
	err := q.Client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: q.CollectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     dimensionEmbedding,
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		return err
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

func (q *QdrantDB) ProcessDocuments(docs []Document, batchSize int) ([]*qdrant.PointStruct, error) {
	var points []*qdrant.PointStruct
	client := ollama.NewClient() //TODO: break this into its own struct/functions

	for i := 0; i < len(docs); i += batchSize {
		end := min(i+batchSize, len(docs))
		batch := docs[i:end]

		for j, doc := range batch {
			chunks := splitIntoChunks(doc.Content, 1000)

			for k, chunk := range chunks {
				embReq := &ollama.EmbeddingRequest{
					Prompt: chunk,
					Model:  "llama3",
				}

				embs, err := client.Embed(context.Background(), embReq)
				if err != nil {
					return nil, fmt.Errorf("embedding error: %w", err)
				}
				embedVals := make([]float32, len(embs[0].Vector)) // Use the Vector field
				for i, val := range embs[0].Vector {
					embedVals[i] = float32(val)
				}
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

func (q *QdrantDB) AddDocuments(points []*qdrant.PointStruct) {
	fmt.Println("AddDocuments")
	operationInfo, err := q.Client.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: q.CollectionName,
		// Points:         []*qdrant.PointStruct{},
		Points: points,
	},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(operationInfo)

}

// func (q *QdrantDB) SearchDocuments(searchText string) ([]*embeddings.Embedding, error) {

// 	client := ollama.NewClient()
// 	embReq := &ollama.EmbeddingRequest{
// 		Prompt: searchText,
// 		Model:  "llama3", // Same model used for document embeddings
// 	}
// 	embs, err := client.Embed(context.Background(), embReq)
// 	if err != nil {
// 		return nil, fmt.Errorf("search embedding error: %w", err)
// 	}

// 	searchVals := make([]float32, len(embs[0].Vector))
// 	for i, val := range embs[0].Vector {
// 		searchVals[i] = float32(val)
// 	}

// 	searchResults, err := q.Client.Query(context.Background(), &qdrant.QueryPoints{
// 		CollectionName: q.CollectionName,
// 		Query:          qdrant.NewQuery(searchVals...),
// 		WithPayload:    &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true}},
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("searching error: %w", err)
// 	}
// 	fmt.Println("searchResults ", searchResults)
// 	for _, searchResult := range searchResults {
// 		fmt.Printf("Found %d matches\n", len(searchResults))
// 		for _, match := range searchResult.Payload {
// 			fmt.Printf("Match (score: %.2f):\n", searchResult.GetScore())
// 			fmt.Printf("Text: %s\n", match.GetListValue().String())
// 		}

// 	}
// 	return nil, nil
// }

func (q *QdrantDB) SearchDocuments(searchText string) ([]*qdrant.ScoredPoint, error) {
	fmt.Println("SearchDocuments")
	client := ollama.NewClient()
	embReq := &ollama.EmbeddingRequest{
		Prompt: searchText,
		Model:  "llama3",
	}
	embs, err := client.Embed(context.Background(), embReq)
	if err != nil {
		return nil, fmt.Errorf("search embedding error: %w", err)
	}

	searchVals := make([]float32, len(embs[0].Vector))
	for i, val := range embs[0].Vector {
		searchVals[i] = float32(val)
	}

	searchResults, err := q.Client.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: q.CollectionName,
		Query:          qdrant.NewQuery(searchVals...),
		WithPayload: &qdrant.WithPayloadSelector{ //explicitly return payload otherwise its empty
			SelectorOptions: &qdrant.WithPayloadSelector_Enable{
				Enable: true,
			}},
	})
	if err != nil {
		return nil, fmt.Errorf("searching error: %w", err)
	}
	// fmt.Println("searchResults ", searchResults)
	fmt.Printf("Found %d matches\n", len(searchResults))
	for _, point := range searchResults {
		fmt.Printf("Score: %.4f\n", point.Score)
		fmt.Printf("ID: %v\n", point.Id)
		fmt.Printf("Payload: %v\n", point.Payload)
		if textValue, ok := point.Payload["text"]; ok {
			// Use the appropriate getter based on how you stored it
			text := textValue.GetStringValue() // or GetKind() to see what type it is
			fmt.Printf("Text: %s\n", text)
		}

		// Get metadata if it exists
		if metaValue, ok := point.Payload["metadata"]; ok {
			meta := metaValue.GetStructValue()
			fmt.Printf("Metadata: %+v\n", meta)
		}

		if point.Vectors != nil {
			x := point.GetVectors()
			fmt.Println("Vector length: ", x)
		}
		fmt.Printf("Version: %d\n", point.Version)
		fmt.Println("---")
	}

	return searchResults, nil
}
