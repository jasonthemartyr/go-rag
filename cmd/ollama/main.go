package main

import (
	"context"
	"fmt"
	"log"

	"github.com/milosgajdos/go-embeddings/ollama"
	"github.com/ollama/ollama/api"
)

func main() {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	x := &api.EmbeddingRequest{
		Model:     "",
		Prompt:    "",
		KeepAlive: &api.Duration{},
		Options:   map[string]interface{}{},
	}
	ctx := context.Background()
	req := &api.EmbedRequest{
		Model: "llama3",
		Input: "what is life",
		// KeepAlive: &api.Duration{},
		// Truncate:  new(bool),
		// Options:   map[string]interface{}{},
	}
	embs, err := client.Embed(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	c := ollama.NewClient()
	embReq := &ollama.EmbeddingRequest{
		Prompt: "what is life",
		Model:  "llama3",
	}
	embs2, err := c.Embed(context.Background(), embReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("EMBS1:")
	fmt.Println(embs.Embeddings)

	fmt.Println("#####")
	fmt.Println("EMBS2:")

	for _, x := range embs2 {
		fmt.Println(x)
	}
	// fmt.Printf("%d", len(embs2))
}

// var (
// 	prompt string
// 	model  string
// )

// func init() {
// 	flag.StringVar(&prompt, "prompt", "what is life", "input prompt")
// 	flag.StringVar(&model, "model", "llama3", "model name")
// }

// func main() {
// 	flag.Parse()

// 	if model == "" {
// 		log.Fatal("missing ollama model")
// 	}

// 	c := ollama.NewClient()

// 	embReq := &ollama.EmbeddingRequest{
// 		Prompt: prompt,
// 		Model:  model,
// 	}

// 	embs, err := c.Embed(context.Background(), embReq)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(embs)
// 	for _, x := range embs {
// 		fmt.Println(x)
// 	}
// 	// fmt.Printf("got %d embeddings", len(embs))
// }
