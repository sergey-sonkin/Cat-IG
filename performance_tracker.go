package main

import (
	"fmt"
	"sort"
	"sync"
)

type PerformanceTracker struct {
	performances []PostPerformance
	mu           sync.RWMutex
}

func NewPerformanceTracker() *PerformanceTracker {
	return &PerformanceTracker{
		performances: make([]PostPerformance, 0),
	}
}

func (pt *PerformanceTracker) AddPerformance(performance PostPerformance) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	pt.performances = append(pt.performances, performance)
}

func (pt *PerformanceTracker) GetBestPerformingPosts(limit int) []PostPerformance {
	if limit <= 0 {
		limit = 10
	}

	pt.mu.RLock()
	defer pt.mu.RUnlock()

	// Create a copy to avoid modifying original slice
	sorted := make([]PostPerformance, len(pt.performances))
	copy(sorted, pt.performances)

	// Sort by engagement rate (descending)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].EngagementRate > sorted[j].EngagementRate
	})

	if len(sorted) < limit {
		return sorted
	}
	return sorted[:limit]
}

func (pt *PerformanceTracker) GetOptimalPostingTimes() []OptimalTimeSlot {
	pt.mu.RLock()
	defer pt.mu.RUnlock()

	timeSlots := make(map[string]*timeSlotData)

	for _, perf := range pt.performances {
		hour := perf.PostedAt.Hour()
		dayOfWeek := int(perf.PostedAt.Weekday())
		key := formatTimeSlotKey(dayOfWeek, hour)

		if slot, exists := timeSlots[key]; exists {
			slot.totalEngagement += perf.EngagementRate
			slot.count++
		} else {
			timeSlots[key] = &timeSlotData{
				hour:            hour,
				dayOfWeek:       dayOfWeek,
				totalEngagement: perf.EngagementRate,
				count:           1,
			}
		}
	}

	// Convert to slice and calculate average performance
	slots := make([]OptimalTimeSlot, 0, len(timeSlots))
	for _, data := range timeSlots {
		slots = append(slots, OptimalTimeSlot{
			Hour:             data.hour,
			DayOfWeek:        data.dayOfWeek,
			PerformanceScore: data.totalEngagement / float64(data.count),
		})
	}

	// Sort by performance score (descending)
	sort.Slice(slots, func(i, j int) bool {
		return slots[i].PerformanceScore > slots[j].PerformanceScore
	})

	// Return top 5
	if len(slots) < 5 {
		return slots
	}
	return slots[:5]
}

func (pt *PerformanceTracker) ShouldPromoteToMain(postPerformance PostPerformance) bool {
	avgEngagement := pt.GetAverageEngagementRate()
	threshold := avgEngagement * 1.5

	return postPerformance.EngagementRate > threshold && postPerformance.Views > 500
}

func (pt *PerformanceTracker) GetAverageEngagementRate() float64 {
	pt.mu.RLock()
	defer pt.mu.RUnlock()

	if len(pt.performances) == 0 {
		return 0.0
	}

	total := 0.0
	for _, perf := range pt.performances {
		total += perf.EngagementRate
	}

	return total / float64(len(pt.performances))
}

func (pt *PerformanceTracker) GetAnalytics() AnalyticsData {
	pt.mu.RLock()
	defer pt.mu.RUnlock()

	totalPosts := len(pt.performances)
	totalViews := 0
	totalLikes := 0

	for _, perf := range pt.performances {
		totalViews += perf.Views
		totalLikes += perf.Likes
	}

	avgEngagement := pt.GetAverageEngagementRate()
	bestPosts := pt.GetBestPerformingPosts(1)
	optimalTimes := pt.GetOptimalPostingTimes()

	var bestPost *PostPerformance
	if len(bestPosts) > 0 {
		bestPost = &bestPosts[0]
	}

	return AnalyticsData{
		TotalPosts:            totalPosts,
		TotalViews:            totalViews,
		TotalLikes:            totalLikes,
		AverageEngagementRate: avgEngagement,
		BestPerformingPost:    bestPost,
		OptimalTimes:          optimalTimes,
	}
}

type AnalyticsData struct {
	TotalPosts            int               `json:"total_posts"`
	TotalViews            int               `json:"total_views"`
	TotalLikes            int               `json:"total_likes"`
	AverageEngagementRate float64           `json:"average_engagement_rate"`
	BestPerformingPost    *PostPerformance  `json:"best_performing_post"`
	OptimalTimes          []OptimalTimeSlot `json:"optimal_times"`
}

type timeSlotData struct {
	hour            int
	dayOfWeek       int
	totalEngagement float64
	count           int
}

func formatTimeSlotKey(dayOfWeek, hour int) string {
	return fmt.Sprintf("%d-%d", dayOfWeek, hour)
}
