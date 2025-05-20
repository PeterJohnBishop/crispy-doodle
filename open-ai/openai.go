package ai

import (
	"context"
	"log"
	"net/http"

	"crispy-doodle/main.go/global"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
)

func OpenAI() *openai.Client {

	client := openai.NewClient(global.OpenAIKey)
	log.Printf("[OPENAI] API Key: %s", global.OpenAIKey)

	log.Printf("[CONNECTED] to OpenAI")
	return client
}

type UserPrompt struct {
	Prompt string `json:"prompt"`
}

func QueryOpenAI(client *openai.Client, c *gin.Context) {
	var input UserPrompt
	if err := c.ShouldBindJSON(&input); err != nil || input.Prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid prompt"})
		return
	}

	req := openai.ChatCompletionRequest{
		Model: openai.GPT4, // or GPT3Dot5Turbo
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: input.Prompt},
		},
		Temperature: 0.7,
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": resp.Choices[0].Message.Content,
	})
}
