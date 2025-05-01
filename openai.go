package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go" // imported as openai
	"github.com/openai/openai-go/option"
)

func Open() *openai.Client {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	key := os.Getenv("OPENAI_API_KEY")

	client := openai.NewClient(
		option.WithAPIKey(key),
	)
	log.Printf("Connected to OpenAI\n")
	return &client
}
