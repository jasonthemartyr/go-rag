package main

import "time"

type (
	ChatMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	ChatRequest struct {
		Model    string        `json:"model"`
		Messages []ChatMessage `json:"messages"`
		Stream   bool          `json:"stream"`
	}
	ChatResponse struct {
		Model              string      `json:"model"`
		CreatedAt          time.Time   `json:"created_at"`
		Message            ChatMessage `json:"message"`
		Done               bool        `json:"done"`
		TotalDuration      int64       `json:"total_duration"`
		LoadDuration       int         `json:"load_duration"`
		PromptEvalCount    int         `json:"prompt_eval_count"`
		PromptEvalDuration int         `json:"prompt_eval_duration"`
		EvalCount          int         `json:"eval_count"`
		EvalDuration       int64       `json:"eval_duration"`
	}
)
