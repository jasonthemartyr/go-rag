# go-rag

Repo meant for deving with RAG and Golang

```bash
docker run -d -p 6333:6333 -p 6334:6334 \
    -v $(pwd)/qdrant_storage:/qdrant/storage:z \
    qdrant/qdrant
```

```bash

sysctl -n hw.ncpu
sysctl hw.memsize | awk '{print $2/1024/1024/1024 " GB"}'

docker run -d \
  --cpus=8 \
  --memory=12g \
  -v ollama:/root/.ollama \
  -p 11434:11434 \
  --name ollama \
  ollama/ollama && docker exec -d ollama ollama run llama3
```

manual query:

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