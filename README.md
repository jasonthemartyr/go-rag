# go-rag

Repo meant for dev'ing with RAG and Golang locally.

## Pre-reqs

- [Ollama container](https://ollama.com/blog/ollama-is-now-available-as-an-official-docker-image)
- [Qdrant Vector DB](https://qdrant.tech/documentation/quickstart/)

Start containers:

```bash
make pre-reqs
```

Test with simple query:

```bash
$ make test-query

{"model":"llama3","created_at":"2024-11-15T19:47:54.75596596Z","message":{"role":"assistant","content":"Hi! It's nice to meet you. Is there something I can help you with, or would you like to chat?"},"done_reason":"stop","done":true,"total_duration":9561533795,"load_duration":39266250,"prompt_eval_count":15,"prompt_eval_duration":4260226000,"eval_count":26,"eval_duration":5217299000}%
```

reference command to optionally set memory/cpu for Ollama:

```bash
sysctl -n hw.ncpu
sysctl hw.memsize | awk '{print $2/1024/1024/1024 " GB"}'
```

Manual query:

```bash

curl localhost:11434/api/chat -d '{
  "model": "llama3",
  "messages": [
    {
      "role": "user",
      "content": "can you say hi?"
    }
  ],
  "stream": false
}'
```

## references

  - https://github.com/golang/example/blob/master/ragserver/ragserver/json.go
  
embeddings:
    - https://github.com/milosgajdos/go-embeddings
    - https://github.com/milosgajdos/go-embeddings/blob/main/ollama/embedding.go

vector DB:
 - https://qdrant.tech/documentation/quickstart/
 - https://github.com/qdrant/go-client
 - https://github.com/milvus-io/milvus-sdk-go
 - https://github.com/PabloSanchi/RAG-GO-Milvus
 - https://zilliz.com/learn