package openai

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go" // imported as openai
	"github.com/openai/openai-go/option"
)

func OpenAI() *openai.Client {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(
		option.WithAPIKey(key),
	)

	log.Printf("[CONNECTED] to OpenAI")
	return &client
}
