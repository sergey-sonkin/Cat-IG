package main

import (
	"context"
	"fmt"
	"log"
	"os"
)

func testVideoGeneration() {
	fmt.Println("🧪 Testing video generation...")

	ctx := context.Background()

	provider := VideoProvider(getEnvWithDefault("VIDEO_PROVIDER", "veo2"))
	fmt.Printf("Using provider: %s\n", provider)

	videoGen, err := NewVideoGenerator(
		os.Getenv("GEMINI_API_KEY"),
		os.Getenv("REPLICATE_API_KEY"),
		os.Getenv("VERTEX_API_KEY"),
		os.Getenv("GOOGLE_PROJECT_ID"),
		provider,
	)
	if err != nil {
		log.Fatalf("Failed to create video generator: %v", err)
	}

	promptGen := NewPromptGenerator(os.Getenv("OPENAI_API_KEY"))

	// Test 1: Generate prompt and video
	fmt.Println("Generating a test prompt...")
	prompt, err := promptGen.GeneratePrompt(ctx)
	if err != nil {
		log.Fatalf("Failed to generate prompt: %v", err)
	}

	fmt.Printf("Generated prompt: \"%s\"\n", prompt.Text)

	fmt.Println("Generating video from prompt...")
	video, err := videoGen.GenerateVideo(ctx, prompt)
	if err != nil {
		log.Fatalf("Failed to generate video: %v", err)
	}

	fmt.Println("✅ Video generation successful!")
	fmt.Printf("Video ID: %s\n", video.ID)
	fmt.Printf("Video URL: %s\n", video.VideoURL)
	fmt.Printf("Duration: %d seconds\n", video.Duration)

	// Test 2: Custom prompt
	fmt.Println("\n🎨 Testing custom prompt...")
	customPrompt := "A sophisticated orange tabby cat wearing tiny glasses, sitting at a miniature desk, looking directly at the camera with an expression that says 'I've seen things you wouldn't believe.' The cat occasionally adjusts its glasses and sighs dramatically."

	fmt.Printf("Testing custom prompt: \"%s\"\n", customPrompt)
	customVideo, err := videoGen.TestGeneration(ctx, customPrompt)
	if err != nil {
		log.Fatalf("Failed to generate custom video: %v", err)
	}

	fmt.Println("✅ Custom prompt test successful!")
	fmt.Printf("Video URL: %s\n", customVideo.VideoURL)

	fmt.Println("\n🎉 All tests completed!")
}

func init() {
	// This allows the file to be run as a standalone test
	if len(os.Args) > 1 && os.Args[1] == "test-video" {
		testVideoGeneration()
		os.Exit(0)
	}
}