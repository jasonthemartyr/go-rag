#!/usr/bin/make -f
makefileDir := $(dir $(lastword $(MAKEFILE_LIST)))
CWD = $(shell cd $(makefileDir) && pwd)

help: ## help using this makefile
	@ grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

pre-reqs: qdrant ollama ## Start both Qdrant and Ollama services
	@echo "Starting all services..."
	@make clean
	@make qdrant
	@make ollama
	@echo "All services started successfully"

qdrant: ## Start Qdrant vector database
	@echo "Starting Qdrant..."
	@docker run -d -p 6333:6333 -p 6334:6334 \
		-v qdrant_storage:/qdrant/storage \
		--name qdrant \
		qdrant/qdrant

ollama: ## Start Ollama with Llama3 model
	@echo "Starting Ollama..."
	@docker run -d --cpus=6 \
		--memory=12g \
		-v ollama:/root/.ollama \
		-p 11434:11434 \
		--name ollama \
		ollama/ollama
	@echo "Waiting for Ollama to start..."
	@sleep 5
	@docker exec -d ollama ollama run llama3

test-query: ## Test Ollama with a simple query
	curl localhost:11434/api/chat -d '{"model":"llama3","messages":[{"role":"user","content":"can you say hi?"}],"stream":false}'

clean: ## Stop and remove containers
	@echo "Cleaning up containers..."
	@docker stop ollama qdrant || exit 0
	@docker rm -f ollama qdrant || exit 0

.PHONY: help pre-reqs qdrant ollama test-query clean