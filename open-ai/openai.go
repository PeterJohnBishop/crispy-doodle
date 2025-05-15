package openai

import (
	"log"

	"crispy-doodle/main.go/global"

	"github.com/openai/openai-go" // imported as openai
	"github.com/openai/openai-go/option"
)

func OpenAI() *openai.Client {

	key := global.OpenAIKey
	client := openai.NewClient(
		option.WithAPIKey(key),
	)

	log.Printf("[CONNECTED] to OpenAI")
	return &client
}
