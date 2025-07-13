package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"github.com/replicate/replicate-go"
	"google.golang.org/api/option"
)

type VideoGenerator struct {
	geminiClient    *genai.Client
	replicateClient *replicate.Client
	vertexAPIKey    string
	projectID       string
	provider        VideoProvider
}

func NewVideoGenerator(geminiAPIKey, replicateAPIKey, vertexAPIKey, projectID string, provider VideoProvider) (*VideoGenerator, error) {
	ctx := context.Background()
	
	geminiClient, err := genai.NewClient(ctx, option.WithAPIKey(geminiAPIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	replicateClient, err := replicate.NewClient(replicate.WithToken(replicateAPIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Replicate client: %w", err)
	}

	return &VideoGenerator{
		geminiClient:    geminiClient,
		replicateClient: replicateClient,
		vertexAPIKey:    vertexAPIKey,
		projectID:       projectID,
		provider:        provider,
	}, nil
}

func (vg *VideoGenerator) GenerateVideo(ctx context.Context, prompt *VideoPrompt) (*GeneratedVideo, error) {
	fmt.Printf("Generating video with %s for prompt: %s\n", vg.provider, prompt.Text)

	switch vg.provider {
	case Veo3Replicate:
		return vg.generateWithVeo3Replicate(ctx, prompt)
	case Veo3Vertex:
		return vg.generateWithVeo3Vertex(ctx, prompt)
	default:
		return vg.generateWithVeo2(ctx, prompt)
	}
}

func (vg *VideoGenerator) generateWithVeo2(ctx context.Context, prompt *VideoPrompt) (*GeneratedVideo, error) {
	model := vg.geminiClient.GenerativeModel("veo-2.0-generate-001")
	
	// Note: This is a placeholder implementation as the actual Gemini video API
	// might have different method signatures in the Go client
	resp, err := model.GenerateContent(ctx, genai.Text(prompt.Text))
	if err != nil {
		return nil, fmt.Errorf("Veo 2 generation failed: %w", err)
	}

	// Extract video URL from response (placeholder logic)
	videoURL := "https://placeholder-video.com/veo2/" + prompt.ID + ".mp4"
	if resp.Candidates != nil && len(resp.Candidates) > 0 {
		// Extract actual video URL from response when available
	}

	return &GeneratedVideo{
		ID:        uuid.New().String(),
		PromptID:  prompt.ID,
		VideoURL:  videoURL,
		Duration:  5,
		CreatedAt: time.Now(),
	}, nil
}

func (vg *VideoGenerator) generateWithVeo3Replicate(ctx context.Context, prompt *VideoPrompt) (*GeneratedVideo, error) {
	input := replicate.PredictionInput{
		"prompt":          prompt.Text,
		"enhance_prompt":  true,
		"negative_prompt": "low quality, blurry, distorted",
	}

	webhook := replicate.Webhook{
		URL:    "",
		Events: []replicate.WebhookEventType{"completed"},
	}

	prediction, err := vg.replicateClient.CreatePrediction(ctx, "google/veo-3", input, &webhook, false)
	if err != nil {
		return nil, fmt.Errorf("Veo 3 Replicate generation failed: %w", err)
	}

	// Wait for completion
	err = vg.replicateClient.Wait(ctx, prediction)
	if err != nil {
		return nil, fmt.Errorf("Veo 3 Replicate wait failed: %w", err)
	}

	videoURL, ok := prediction.Output.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected output format from Veo 3 Replicate")
	}

	return &GeneratedVideo{
		ID:        uuid.New().String(),
		PromptID:  prompt.ID,
		VideoURL:  videoURL,
		Duration:  8,
		CreatedAt: time.Now(),
	}, nil
}

func (vg *VideoGenerator) generateWithVeo3Vertex(ctx context.Context, prompt *VideoPrompt) (*GeneratedVideo, error) {
	endpoint := fmt.Sprintf("https://aiplatform.googleapis.com/v1/projects/%s/locations/us-central1/publishers/google/models/veo-3.0-generate-preview:predictLongRunning", vg.projectID)

	requestBody := map[string]interface{}{
		"instances": []map[string]interface{}{
			{"prompt": prompt.Text},
		},
		"parameters": map[string]interface{}{
			"aspectRatio":     "9:16",
			"durationSeconds": 8,
			"numberOfVideos":  1,
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+vg.vertexAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Vertex AI request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Vertex AI request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	operationName, ok := result["name"].(string)
	if !ok {
		return nil, fmt.Errorf("no operation name in response")
	}

	videoURL, err := vg.pollVertexOperation(ctx, operationName)
	if err != nil {
		return nil, fmt.Errorf("failed to poll operation: %w", err)
	}

	return &GeneratedVideo{
		ID:        uuid.New().String(),
		PromptID:  prompt.ID,
		VideoURL:  videoURL,
		Duration:  8,
		CreatedAt: time.Now(),
	}, nil
}

func (vg *VideoGenerator) pollVertexOperation(ctx context.Context, operationName string) (string, error) {
	maxAttempts := 30
	delay := 10 * time.Second

	for attempt := 0; attempt < maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}

		endpoint := fmt.Sprintf("https://aiplatform.googleapis.com/v1/%s", operationName)
		req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
		if err != nil {
			return "", fmt.Errorf("failed to create poll request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+vg.vertexAPIKey)

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("failed to poll operation: %w", err)
		}

		var operation map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&operation); err != nil {
			resp.Body.Close()
			return "", fmt.Errorf("failed to decode poll response: %w", err)
		}
		resp.Body.Close()

		if done, ok := operation["done"].(bool); ok && done {
			if errorObj, exists := operation["error"]; exists {
				return "", fmt.Errorf("operation failed: %v", errorObj)
			}

			if response, ok := operation["response"].(map[string]interface{}); ok {
				if predictions, ok := response["predictions"].([]interface{}); ok && len(predictions) > 0 {
					if prediction, ok := predictions[0].(map[string]interface{}); ok {
						if videoURL, ok := prediction["videoUrl"].(string); ok {
							return videoURL, nil
						}
					}
				}
			}

			return "", fmt.Errorf("no video URL in completed operation")
		}

		fmt.Printf("Video generation in progress... (attempt %d/%d)\n", attempt+1, maxAttempts)
		time.Sleep(delay)
	}

	return "", fmt.Errorf("video generation timed out")
}

func (vg *VideoGenerator) GenerateBatch(ctx context.Context, prompts []*VideoPrompt) ([]*GeneratedVideo, error) {
	videos := make([]*GeneratedVideo, 0, len(prompts))
	
	for _, prompt := range prompts {
		video, err := vg.GenerateVideo(ctx, prompt)
		if err != nil {
			fmt.Printf("Failed to generate video for prompt %s: %v\n", prompt.ID, err)
			continue
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func (vg *VideoGenerator) SetProvider(provider VideoProvider) {
	vg.provider = provider
	fmt.Printf("Switched to %s for video generation\n", provider)
}

func (vg *VideoGenerator) TestGeneration(ctx context.Context, testPrompt string) (*GeneratedVideo, error) {
	if testPrompt == "" {
		testPrompt = "A fluffy orange cat wearing sunglasses sits in a director's chair, occasionally looking directly at the camera with an expression of existential contemplation"
	}

	prompt := &VideoPrompt{
		ID:        uuid.New().String(),
		Text:      testPrompt,
		Theme:     "test",
		CreatedAt: time.Now(),
	}

	return vg.GenerateVideo(ctx, prompt)
}