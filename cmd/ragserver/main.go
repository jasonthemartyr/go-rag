package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const defaultOllamaURL = "http://localhost:11434/api/chat"

// docker run -d -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama && docker exec -d ollama ollama run llama3

func main() {
	start := time.Now()
	msg := ChatMessage{
		Role:    "user",
		Content: "Why is the sky blue?",
	}
	req := ChatRequest{
		Model:    "llama3.1",
		Stream:   false,
		Messages: []ChatMessage{msg},
	}
	resp, err := talkToOllama(defaultOllamaURL, req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(resp.Message.Content)
	fmt.Printf("Completed in %v", time.Since(start))
}

func talkToOllama(url string, ollamaReq ChatRequest) (*ChatResponse, error) {
	js, err := json.Marshal(&ollamaReq)
	if err != nil {
		return nil, err
	}
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
