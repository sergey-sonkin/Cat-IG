package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	fmt.Println("🧪 AI Cat Video Generation Test")
	fmt.Println("================================")

	ctx := context.Background()

	// Check required environment variables
	if err := checkRequiredEnvVars(); err != nil {
		log.Fatalf("❌ Environment setup error: %v", err)
	}

	provider := VideoProvider(getEnvWithDefault("VIDEO_PROVIDER", "veo2"))
	fmt.Printf("📹 Using video provider: %s\n\n", provider)

	// Initialize components
	fmt.Println("🔧 Initializing components...")

	promptGen := NewPromptGenerator(os.Getenv("OPENAI_API_KEY"))
	fmt.Println("✅ Prompt generator ready")

	videoGen, err := NewVideoGenerator(
		os.Getenv("GEMINI_API_KEY"),
		os.Getenv("REPLICATE_API_KEY"),
		os.Getenv("VERTEX_API_KEY"),
		os.Getenv("GOOGLE_PROJECT_ID"),
		provider,
	)
	if err != nil {
		log.Fatalf("❌ Failed to create video generator: %v", err)
	}
	fmt.Println("✅ Video generator ready")

	// Test 1: Generate some prompts
	fmt.Println("\n🎭 Generating test prompts...")
	prompts, err := promptGen.GenerateBatch(ctx, 3)
	if err != nil {
		log.Fatalf("❌ Failed to generate prompts: %v", err)
	}

	fmt.Printf("✅ Generated %d prompts:\n", len(prompts))
	for i, prompt := range prompts {
		fmt.Printf("   %d. [%s] %s\n", i+1, prompt.Theme, prompt.Text)
	}

	// Test 2: Generate a video from the first prompt
	if len(prompts) > 0 {
		fmt.Printf("\n🎬 Generating video from first prompt...\n")
		fmt.Printf("Prompt: \"%s\"\n", prompts[0].Text)

		startTime := time.Now()
		video, err := videoGen.GenerateVideo(ctx, prompts[0])
		duration := time.Since(startTime)

		if err != nil {
			fmt.Printf("❌ Video generation failed: %v\n", err)
		} else {
			fmt.Printf("✅ Video generated successfully!\n")
			fmt.Printf("   📁 Video ID: %s\n", video.ID)
			fmt.Printf("   🔗 Video URL: %s\n", video.VideoURL)
			fmt.Printf("   ⏱️  Duration: %d seconds\n", video.Duration)
			fmt.Printf("   🕐 Generation time: %v\n", duration)
		}
	}

	// Test 3: Custom prompt test
	fmt.Println("\n🎨 Testing custom prompt...")
	customPrompt := "A sophisticated orange tabby cat wearing tiny reading glasses sits at a miniature wooden desk with a tiny laptop. The cat types furiously with one paw while occasionally glancing at the camera with a look that says 'I'm handling very important cat business here.' The cat sighs dramatically and adjusts its glasses."

	fmt.Printf("Custom prompt: \"%s\"\n", customPrompt)

	startTime := time.Now()
	customVideo, test_gen_err := videoGen.TestGeneration(ctx, customPrompt)
	duration := time.Since(startTime)

	if test_gen_err != nil {
		fmt.Printf("❌ Custom video generation failed: %v\n", test_gen_err)
	} else {
		fmt.Printf("✅ Custom video generated successfully!\n")
		fmt.Printf("   📁 Video ID: %s\n", customVideo.ID)
		fmt.Printf("   🔗 Video URL: %s\n", customVideo.VideoURL)
		fmt.Printf("   ⏱️  Duration: %d seconds\n", customVideo.Duration)
		fmt.Printf("   🕐 Generation time: %v\n", duration)
	}

	fmt.Println("\n🎉 Video generation tests complete!")
	fmt.Println("\nNext steps:")
	fmt.Println("• Check the video URLs above to see your generated content")
	fmt.Println("• Try different VIDEO_PROVIDER values in .env (veo2, veo3-replicate, veo3-vertex)")
	fmt.Println("• Run 'go run main_full.go' to test the complete Instagram posting pipeline")
}

func checkRequiredEnvVars() error {
	required := map[string]string{
		"OPENAI_API_KEY": "OpenAI API key for prompt generation",
	}

	provider := getEnvWithDefault("VIDEO_PROVIDER", "veo2")

	switch VideoProvider(provider) {
	case Veo2:
		required["GEMINI_API_KEY"] = "Gemini API key for Veo 2 video generation"
	case Veo3Replicate:
		required["REPLICATE_API_KEY"] = "Replicate API key for Veo 3 video generation"
	case Veo3Vertex:
		required["VERTEX_API_KEY"] = "Vertex AI access token for Veo 3 video generation"
		required["GOOGLE_PROJECT_ID"] = "Google Cloud project ID for Vertex AI"
	}

	var missing []string
	for env, description := range required {
		if os.Getenv(env) == "" {
			missing = append(missing, fmt.Sprintf("  %s: %s", env, description))
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables:\n%s\n\nPlease set these in your .env file",
			fmt.Sprintf("%s", missing))
	}

	return nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
