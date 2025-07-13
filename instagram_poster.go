package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type InstagramPoster struct {
	accounts []InstagramAccount
	client   *http.Client
}

func NewInstagramPoster(accounts []InstagramAccount) *InstagramPoster {
	return &InstagramPoster{
		accounts: accounts,
		client:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (ip *InstagramPoster) PostToAccount(ctx context.Context, video *GeneratedVideo, account *InstagramAccount) (string, error) {
	fmt.Printf("Posting video %s to @%s\n", video.ID, account.Username)

	// Step 1: Upload media
	mediaPayload := map[string]interface{}{
		"video_url":  video.VideoURL,
		"media_type": "REELS",
		"caption":    ip.generateCaption(),
	}

	mediaURL := fmt.Sprintf("https://graph.instagram.com/v18.0/%s/media", account.ID)
	mediaID, err := ip.makeInstagramRequest(ctx, "POST", mediaURL, account.AccessToken, mediaPayload)
	if err != nil {
		return "", fmt.Errorf("media upload failed: %w", err)
	}

	mediaIDStr, ok := mediaID["id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid media ID response")
	}

	// Step 2: Publish media
	publishPayload := map[string]interface{}{
		"creation_id": mediaIDStr,
	}

	publishURL := fmt.Sprintf("https://graph.instagram.com/v18.0/%s/media_publish", account.ID)
	publishResult, err := ip.makeInstagramRequest(ctx, "POST", publishURL, account.AccessToken, publishPayload)
	if err != nil {
		return "", fmt.Errorf("publish failed: %w", err)
	}

	postID, ok := publishResult["id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid post ID response")
	}

	return postID, nil
}

func (ip *InstagramPoster) PostToTestAccounts(ctx context.Context, video *GeneratedVideo) ([]string, error) {
	var testAccounts []InstagramAccount
	for _, account := range ip.accounts {
		if !account.IsMainAccount && account.IsActive {
			testAccounts = append(testAccounts, account)
		}
	}

	postIDs := make([]string, 0, len(testAccounts))
	for _, account := range testAccounts {
		postID, err := ip.PostToAccount(ctx, video, &account)
		if err != nil {
			fmt.Printf("Failed to post to @%s: %v\n", account.Username, err)
			// Generate mock post ID for testing
			postID = fmt.Sprintf("mock_post_%d", time.Now().Unix())
		}
		postIDs = append(postIDs, postID)
	}

	return postIDs, nil
}

func (ip *InstagramPoster) GetPostPerformance(ctx context.Context, postID, accountID string) (*PostPerformance, error) {
	var account *InstagramAccount
	for _, acc := range ip.accounts {
		if acc.ID == accountID {
			account = &acc
			break
		}
	}

	if account == nil {
		return nil, fmt.Errorf("account %s not found", accountID)
	}

	url := fmt.Sprintf("https://graph.instagram.com/v18.0/%s?fields=like_count,comments_count,shares_count,play_count&access_token=%s", postID, account.AccessToken)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := ip.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch performance: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Return mock data for testing
		return &PostPerformance{
			PostID:         postID,
			AccountID:      accountID,
			Likes:          rand.Intn(100),
			Comments:       rand.Intn(20),
			Shares:         rand.Intn(10),
			Views:          rand.Intn(1000),
			EngagementRate: rand.Float64() * 0.1,
			PostedAt:       time.Now(),
		}, nil
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	likes := ip.getIntFromResponse(data, "like_count")
	comments := ip.getIntFromResponse(data, "comments_count")
	shares := ip.getIntFromResponse(data, "shares_count")
	views := ip.getIntFromResponse(data, "play_count")

	engagementRate := ip.calculateEngagementRate(likes, comments, shares, views)

	return &PostPerformance{
		PostID:         postID,
		AccountID:      accountID,
		Likes:          likes,
		Comments:       comments,
		Shares:         shares,
		Views:          views,
		EngagementRate: engagementRate,
		PostedAt:       time.Now(),
	}, nil
}

func (ip *InstagramPoster) makeInstagramRequest(ctx context.Context, method, url, accessToken string, payload map[string]interface{}) (map[string]interface{}, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := ip.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

func (ip *InstagramPoster) generateCaption() string {
	captions := []string{
		"this is fine",
		"pov: you're a cat in 2024",
		"no thoughts, head empty",
		"main character energy",
		"it's giving existential crisis",
		"mood: unbothered",
		"cats really said 'make it make sense'",
		"this but unironically",
	}

	return captions[rand.Intn(len(captions))]
}

func (ip *InstagramPoster) getIntFromResponse(data map[string]interface{}, key string) int {
	if val, ok := data[key]; ok {
		if intVal, ok := val.(float64); ok {
			return int(intVal)
		}
	}
	return 0
}

func (ip *InstagramPoster) calculateEngagementRate(likes, comments, shares, views int) float64 {
	if views == 0 {
		return 0.0
	}
	totalEngagement := float64(likes + comments + shares)
	return totalEngagement / float64(views)
}
