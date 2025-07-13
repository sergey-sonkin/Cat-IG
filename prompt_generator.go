package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	openai "github.com/sashabaranov/go-openai"
)

type PromptGenerator struct {
	client *openai.Client
	themes []string
	situations []string
}

func NewPromptGenerator(apiKey string) *PromptGenerator {
	return &PromptGenerator{
		client: openai.NewClient(apiKey),
		themes: []string{
			"existential dread",
			"corporate middle management",
			"gen z slang misuse",
			"cryptocurrency obsession",
			"wellness influencer parody",
			"linkedin motivational posts",
			"artisanal everything",
			"sustainable living anxiety",
			"dating app failures",
			"work from home chaos",
		},
		situations: []string{
			"realizes they've been eating the same cardboard for 3 years",
			"discovers their humans are just large, hairless cats",
			"starts a podcast about the futility of chasing laser dots",
			"becomes a life coach for other cats",
			"opens a meditation retreat for anxious house pets",
			"launches a startup selling cardboard boxes as premium furniture",
			"writes passive-aggressive emails to their food dispenser",
			"practices mindfulness while knocking things off tables",
			"develops an elaborate conspiracy theory about vacuum cleaners",
			"starts a support group for cats with imposter syndrome",
		},
	}
}

func (pg *PromptGenerator) GeneratePrompt(ctx context.Context) (*VideoPrompt, error) {
	theme := pg.themes[rand.Intn(len(pg.themes))]
	situation := pg.situations[rand.Intn(len(pg.situations))]

	systemPrompt := "You are a creative director for post-ironic cat content. Generate absurd, slightly meta video prompts that combine internet culture with cat behavior. Keep it weird but family-friendly."
	userPrompt := fmt.Sprintf("Create a short video prompt (1-2 sentences) about a cat dealing with \"%s\" where the cat %s. Make it absurd and slightly self-aware.", theme, situation)

	resp, err := pg.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: userPrompt},
		},
		MaxTokens:   150,
		Temperature: 0.9,
	})

	if err != nil {
		fmt.Printf("Failed to generate prompt: %v\n", err)
		// Fallback prompt
		return &VideoPrompt{
			ID:        uuid.New().String(),
			Text:      fmt.Sprintf("A cat %s while contemplating %s, occasionally making direct eye contact with the camera to break the fourth wall.", situation, theme),
			Theme:     theme,
			CreatedAt: time.Now(),
		}, nil
	}

	promptText := resp.Choices[0].Message.Content
	if promptText == "" {
		promptText = "A cat stares judgmentally at the camera while questioning the meaning of existence."
	}

	return &VideoPrompt{
		ID:        uuid.New().String(),
		Text:      promptText,
		Theme:     theme,
		CreatedAt: time.Now(),
	}, nil
}

func (pg *PromptGenerator) GenerateBatch(ctx context.Context, count int) ([]*VideoPrompt, error) {
	prompts := make([]*VideoPrompt, 0, count)
	
	for i := 0; i < count; i++ {
		prompt, err := pg.GeneratePrompt(ctx)
		if err != nil {
			fmt.Printf("Failed to generate prompt %d: %v\n", i+1, err)
			continue
		}
		prompts = append(prompts, prompt)
	}

	return prompts, nil
}