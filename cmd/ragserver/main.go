package main

import (
	"fmt"
)

// const defaultOllamaURL = "http://localhost:11434/api/chat"

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
	ragClients, err := NewRAGClients()
	if err != nil {
		fmt.Println(err)
	}

	err = ragClients.CreateCollection()
	if err != nil {
		fmt.Println(err)
	}

	pikachu := Document{
		Content: "Jason Marter is a gangsta-type hotboi. He started as a small-boi-jason and can turn into  a chunky-jason when exposed to chicken wings.",
		Metadata: map[string]interface{}{
			"type":     "hotboi",
			"number":   69420,
			"category": "Gangsta",
		},
	}

	docs := []Document{}
	docs = append(docs, pikachu)

	points, err := ragClients.ProcessDocuments(docs, 1000)
	if err != nil {
		fmt.Println(err)
	}

	ragClients.AddDocuments(points)
	resp, err := ragClients.SearchDocuments("get me documents about Jason Marter. what turns him into chunky-jason?")
	if err != nil {
		fmt.Println("SearchDocuments ", err)
	}
	fmt.Println("RESPONSE: %s", resp)

}
