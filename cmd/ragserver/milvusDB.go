package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

const (
	milvusHost = "localhost:19530"
	// collectionName = "documents"
	// dimensionEmbedding = 384
	// dimensionEmbedding = 128
)

type MilvusDB struct {
	Client  client.Client
	Context context.Context
}

func NewMilvusDB() (*MilvusDB, error) {
	fmt.Println("create new DB")
	ctx := context.Background()
	client, err := client.NewGrpcClient(ctx, milvusHost)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Milvus: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return &MilvusDB{
		Client:  client,
		Context: ctx,
	}, nil

}

// https://github.com/milvus-io/milvus-sdk-go/blob/master/examples/basic/basic.go
// collection is like a table or container the groups related vectors/content
func (m *MilvusDB) CreateCollection() error {
	fmt.Println("create new collection")
	ctx := m.Context
	collectionName := collectionName
	dim := dimensionEmbedding
	c := m.Client
	collExists, err := c.HasCollection(ctx, collectionName)
	if err != nil {
		return fmt.Errorf("failed to check collection exists: %w", err)
	}
	if collExists {
		// let's say the example collection is only for sampling the API
		// drop old one in case early crash or something
		_ = c.DropCollection(ctx, collectionName)
	}

	schema := &entity.Schema{
		CollectionName: collectionName,
		Description:    "Document store for RAG",
		Fields: []*entity.Field{
			{
				Name:       "id",
				DataType:   entity.FieldTypeInt64,
				PrimaryKey: true,
				AutoID:     true,
			},
			{
				Name:     "content",
				DataType: entity.FieldTypeVarChar,
				// MaxLength: 65535,
			},
			{
				Name:     "embedding",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": fmt.Sprintf("%d", dim),
				},
			},
		},
	}

	err = c.CreateCollection(ctx, schema, entity.DefaultShardNumber)
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}
	return nil

}

// Index = data struct that organize vectors/data for searching
func (m *MilvusDB) CreateIndex() error {
	fmt.Println("create new index")
	ctx := m.Context
	c := m.Client
	idx, err := entity.NewIndexIvfFlat(entity.L2, 1024)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	err = c.CreateIndex(ctx, collectionName, "embedding", idx, false)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	return nil

}

func (m *MilvusDB) InsertData(content []string) error {
	fmt.Println("inserting data")
	ctx := m.Context
	c := m.Client
	embeddings := [][]float32{{0.1, 0.2 /* ... up to dim values */}}
	err := c.LoadCollection(ctx, collectionName, false)
	if err != nil {
		log.Fatal("Failed to load collection:", err.Error())
	}

	//TODO: figure out exactly what this means instead of just copying from an example lol
	contentColumn := entity.NewColumnVarChar("content", content)
	embeddingColumn := entity.NewColumnFloatVector("embedding", len(embeddings[0]), embeddings)
	_, err = c.Insert(ctx, collectionName, "", contentColumn, embeddingColumn)
	if err != nil {
		log.Fatal("Failed to insert:", err.Error())
	}

	return nil
}

func (m *MilvusDB) SearchData(content []string) error {
	fmt.Println("searching data")
	ctx := m.Context
	c := m.Client
	embedding := []float32{0.1, 0.2, 0.3, 0.4}
	searchVec := entity.FloatVector(embedding)
	sp, err := entity.NewIndexIvfFlatSearchParam(10) // nprobe = 10
	if err != nil {
		log.Fatal("Failed to create search params:", err.Error())
	}
	searchResult, err := c.Search(
		ctx,                        // context
		collectionName,             // collection name
		[]string{"content"},        // output fields
		"",                         // partition name
		[]string{},                 // expr
		[]entity.Vector{searchVec}, // vectors to search
		"embedding",                // vector field
		entity.L2,                  // metric type
		2,                          // topK
		sp,                         // search params
	)
	if err != nil {
		log.Fatal("Failed to search:", err.Error())
	}
	for _, result := range searchResult {
		for _, field := range result.Fields {
			// if field.Name() == "content"{

			// }
			fmt.Println(field.Name())
			fmt.Println(field.FieldData().GetVectors().GetSparseFloatVector().GetContents())

		}
	}

	return nil
}
