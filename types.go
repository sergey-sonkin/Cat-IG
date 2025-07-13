package main

import "time"

type VideoPrompt struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Theme     string    `json:"theme"`
	CreatedAt time.Time `json:"created_at"`
}

type GeneratedVideo struct {
	ID        string    `json:"id"`
	PromptID  string    `json:"prompt_id"`
	VideoURL  string    `json:"video_url"`
	Duration  int       `json:"duration"`
	CreatedAt time.Time `json:"created_at"`
}

type InstagramAccount struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	AccessToken   string `json:"access_token"`
	IsMainAccount bool   `json:"is_main_account"`
	IsActive      bool   `json:"is_active"`
}

type PostPerformance struct {
	PostID         string    `json:"post_id"`
	AccountID      string    `json:"account_id"`
	Likes          int       `json:"likes"`
	Comments       int       `json:"comments"`
	Shares         int       `json:"shares"`
	Views          int       `json:"views"`
	EngagementRate float64   `json:"engagement_rate"`
	PostedAt       time.Time `json:"posted_at"`
}

type OptimalTimeSlot struct {
	Hour             int     `json:"hour"`
	DayOfWeek        int     `json:"day_of_week"`
	PerformanceScore float64 `json:"performance_score"`
}

type VideoProvider string

const (
	Veo2           VideoProvider = "veo2"
	Veo3Replicate  VideoProvider = "veo3-replicate"
	Veo3Vertex     VideoProvider = "veo3-vertex"
)