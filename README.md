# go-rag

Repo meant for deving with RAG and Golang

```bash
docker run -d -p 6333:6333 -p 6334:6334 \
    -v $(pwd)/qdrant_storage:/qdrant/storage:z \
    qdrant/qdrant
```


```bash
docker run -d --name milvus_standalone -p 19530:19530 -p 9091:9091 milvusdb/milvus:v2.3.3 milvus run standalone


docker run -d --name milvus_standalone \
  -p 19530:19530 \
  -p 9091:9091 \
  -v $(pwd)/milvus_data:/var/lib/milvus \
  milvusdb/milvus:v2.3.3 milvus run standalone
```

## references

  - https://github.com/milvus-io/milvus-sdk-go
  - https://github.com/PabloSanchi/RAG-GO-Milvus
  - https://github.com/golang/example/blob/master/ragserver/ragserver/json.go
  - https://github.com/qdrant/go-client

embeddings:
https://github.com/milosgajdos/go-embeddings/blob/main/ollama/embedding.go