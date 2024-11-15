package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// const defaultOllamaURL = "http://localhost:11434/api/chat"

// docker run -d -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama && docker exec -d ollama ollama run llama3

func main() {
	ragClients, err := NewRAGClients()
	if err != nil {
		fmt.Println(err)
	}

	err = ragClients.CreateCollection()
	if err != nil {
		fmt.Println(err)
	}

	jsonBoi := Document{
		Content: "Jason is a gangsta-type hotboi. He started as a small-boi-jason and can turn into  a chunky-jason when exposed to chicken wings.",
		Metadata: map[string]interface{}{
			"type":     "hotboi",
			"number":   69420,
			"category": "Gangsta",
		},
	}

	docs := []Document{}
	docs = append(docs, jsonBoi)

	points, err := ragClients.ProcessDocuments(docs, 1000)
	if err != nil {
		fmt.Println(err)
	}

	ragClients.AddDocuments(points)
	// resp, err := ragClients.SearchDocuments("get me documents about Jason. what turns him into chunky-jason?")
	// if err != nil {
	// 	fmt.Println("SearchDocuments ", err)
	// }
	// fmt.Println("RESPONSE: %s", resp)

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\nType your questions and press Enter. Type 'exit' to quit.\n")

	for {
		fmt.Print("Question > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}
		input = strings.TrimSpace(input)
		if strings.ToLower(input) == "exit" {
			fmt.Println("bye I guess...")
			break
		}
		if input == "" {
			continue
		}
		resp, err := ragClients.SearchDocuments(input)
		if err != nil {
			fmt.Println("SearchDocuments ", err)
		}
		fmt.Printf("\nRESPONSE:\n%s\n\n", resp)

	}

}
