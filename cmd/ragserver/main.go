package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const defaultOllamaURL = "http://localhost:11434/api/chat"

// docker run -d -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama && docker exec -d ollama ollama run llama3

func main() {
	// start := time.Now()
	// msg := ChatMessage{
	// 	Role:    "user",
	// 	Content: "Why is the sky blue?",
	// }
	// req := ChatRequest{
	// 	Model:    "llama3",
	// 	Stream:   false,
	// 	Messages: []ChatMessage{msg},
	// }
	// resp, err := talkToOllama(defaultOllamaURL, req)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Println("here")
	// fmt.Println("msg: ", resp.Message.Content)
	// fmt.Printf("Completed in %v", time.Since(start))

	// mDB, err := NewMilvusDB()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// err = mDB.CreateCollection()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// err = mDB.CreateIndex()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// err = mDB.InsertData([]string{"this is a test"})
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = mDB.SearchData([]string{"this is a test"})
	// if err != nil {
	// 	fmt.Println(err)
	// }
	vDB, err := NewQdrantDB()
	if err != nil {
		fmt.Println(err)
	}

	err = vDB.CreateCollection()

	pikachu := Document{
		Content: "Pikachu is an Electric-type Pokémon introduced in Generation I. It evolves from Pichu when leveled up with high friendship and evolves into Raichu when exposed to a Thunder Stone.",
		Metadata: map[string]interface{}{
			"type":       "pokemon",
			"generation": 1,
			"number":     25,
			"category":   "Mouse Pokémon",
		},
	}

	docs := []Document{}
	docs = append(docs, pikachu)
	points, err := vDB.ProcessDocuments(docs, 100)
	if err != nil {
		fmt.Println(err)
	}

	vDB.AddDocuments(points)
	// fmt.Println(x)
	_, err = vDB.SearchDocuments("find me documents about Pikachu")
	if err != nil {
		fmt.Println("SearchDocuments ", err)
	}

	// fmt.Println(docs)

}

func talkToOllama(url string, ollamaReq ChatRequest) (*ChatResponse, error) {
	js, err := json.Marshal(&ollamaReq)
	if err != nil {
		return nil, err
	}
	fmt.Println(ollamaReq)
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
