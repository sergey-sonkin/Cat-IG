package main

import (
	"context"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("🐱 AI Cat Content Generator starting...")

	ctx := context.Background()

	// Initialize components
	promptGen := NewPromptGenerator(os.Getenv("OPENAI_API_KEY"))
	
	videoGen, err := NewVideoGenerator(
		os.Getenv("GEMINI_API_KEY"),
		os.Getenv("REPLICATE_API_KEY"),
		os.Getenv("VERTEX_API_KEY"),
		os.Getenv("GOOGLE_PROJECT_ID"),
		VideoProvider(getEnvWithDefault("VIDEO_PROVIDER", "veo2")),
	)
	if err != nil {
		log.Fatalf("Failed to create video generator: %v", err)
	}

	tracker := NewPerformanceTracker()

	testAccounts := []InstagramAccount{
		{
			ID:            "test1",
			Username:      "cat_vibes_1",
			AccessToken:   os.Getenv("INSTA_TOKEN_1"),
			IsMainAccount: false,
			IsActive:      true,
		},
		{
			ID:            "test2",
			Username:      "cat_vibes_2",
			AccessToken:   os.Getenv("INSTA_TOKEN_2"),
			IsMainAccount: false,
			IsActive:      true,
		},
		{
			ID:            "main",
			Username:      "main_cat_account",
			AccessToken:   os.Getenv("INSTA_TOKEN_MAIN"),
			IsMainAccount: true,
			IsActive:      true,
		},
	}

	poster := NewInstagramPoster(testAccounts)

	// Generate content
	prompts, err := promptGen.GenerateBatch(ctx, 3)
	if err != nil {
		log.Fatalf("Failed to generate prompts: %v", err)
	}

	fmt.Printf("Generated %d prompts:\n", len(prompts))
	for i, prompt := range prompts {
		fmt.Printf("%d. %s\n", i+1, prompt.Text)
	}

	videos, err := videoGen.GenerateBatch(ctx, prompts)
	if err != nil {
		log.Fatalf("Failed to generate videos: %v", err)
	}

	fmt.Printf("Generated %d videos\n", len(videos))

	// Post to test accounts
	for _, video := range videos {
		postIDs, err := poster.PostToTestAccounts(ctx, video)
		if err != nil {
			fmt.Printf("Warning: Failed to post video %s: %v\n", video.ID, err)
			continue
		}
		fmt.Printf("Posted video %s to %d test accounts\n", video.ID, len(postIDs))
	}

	fmt.Println("✅ Content generation and posting complete!")
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}