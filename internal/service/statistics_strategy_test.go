package service

import (
	"testing"

	"social-media-analyzer/internal/models"
)

// TestAggregateStatsStrategy tests aggregate statistics calculation
func TestAggregateStatsStrategy(t *testing.T) {
	strategy := &AggregateStatsStrategy{}

	tests := []struct {
		name              string
		posts             []models.Post
		expectedTotalLikes int
		expectedAvgLikes   float64
		expectedMaxLikes   int
	}{
		{
			name:              "Empty posts",
			posts:             []models.Post{},
			expectedTotalLikes: 0,
			expectedAvgLikes:   0,
			expectedMaxLikes:   0,
		},
		{
			name: "Single post",
			posts: []models.Post{
				{Likes: 100, Comments: 10},
			},
			expectedTotalLikes: 100,
			expectedAvgLikes:   100.0,
			expectedMaxLikes:   100,
		},
		{
			name: "Multiple posts",
			posts: []models.Post{
				{Likes: 100, Comments: 10},
				{Likes: 200, Comments: 20},
				{Likes: 300, Comments: 30},
			},
			expectedTotalLikes: 600,
			expectedAvgLikes:   200.0,
			expectedMaxLikes:   300,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strategy.Calculate(tt.posts)
			stats := result.(AggregateStats)

			if stats.TotalLikes != tt.expectedTotalLikes {
				t.Errorf("Expected total likes %d, got %d", tt.expectedTotalLikes, stats.TotalLikes)
			}
			if stats.AvgLikesPerPost != tt.expectedAvgLikes {
				t.Errorf("Expected avg likes %.1f, got %.1f", tt.expectedAvgLikes, stats.AvgLikesPerPost)
			}
			if stats.MaxLikesPerPost != tt.expectedMaxLikes {
				t.Errorf("Expected max likes %d, got %d", tt.expectedMaxLikes, stats.MaxLikesPerPost)
			}
		})
	}
}

// TestEngagementRateStrategy tests engagement rate calculation
func TestEngagementRateStrategy(t *testing.T) {
	strategy := &EngagementRateStrategy{}

	tests := []struct {
		name               string
		posts              []models.Post
		expectedEngagement int
		expectedAvgViews   int
	}{
		{
			name:               "Empty posts",
			posts:              []models.Post{},
			expectedEngagement: 0,
			expectedAvgViews:   0,
		},
		{
			name: "High engagement posts",
			posts: []models.Post{
				{Likes: 100, Comments: 20, Views: 1000, Reactions: 50},
				{Likes: 150, Comments: 30, Views: 1500, Reactions: 75},
			},
			expectedEngagement: 425, // (100+20+50) + (150+30+75) = 425
			expectedAvgViews:   1250, // (1000 + 1500) / 2
		},
		{
			name: "Low engagement posts",
			posts: []models.Post{
				{Likes: 10, Comments: 2, Views: 100, Reactions: 5},
			},
			expectedEngagement: 17, // 10 + 2 + 5
			expectedAvgViews:   100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strategy.Calculate(tt.posts)
			rate := result.(EngagementRate)

			if rate.TotalEngagement != tt.expectedEngagement {
				t.Errorf("Expected engagement %d, got %d", tt.expectedEngagement, rate.TotalEngagement)
			}
			if rate.AvgViews != tt.expectedAvgViews {
				t.Errorf("Expected avg views %d, got %d", tt.expectedAvgViews, rate.AvgViews)
			}
		})
	}
}

// TestEngagementRateStrategy_LikeToCommentRatio tests like-to-comment ratio calculation
func TestEngagementRateStrategy_LikeToCommentRatio(t *testing.T) {
	strategy := &EngagementRateStrategy{}

	tests := []struct {
		name          string
		posts         []models.Post
		expectedRatio float64
	}{
		{
			name:          "No comments",
			posts:         []models.Post{{Likes: 100, Comments: 0}},
			expectedRatio: 0.0,
		},
		{
			name:          "Equal likes and comments",
			posts:         []models.Post{{Likes: 100, Comments: 100}},
			expectedRatio: 1.0,
		},
		{
			name:          "More likes than comments",
			posts:         []models.Post{{Likes: 200, Comments: 50}},
			expectedRatio: 4.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strategy.Calculate(tt.posts)
			rate := result.(EngagementRate)

			if rate.LikeToCommentRatio != tt.expectedRatio {
				t.Errorf("Expected ratio %.1f, got %.1f", tt.expectedRatio, rate.LikeToCommentRatio)
			}
		})
	}
}

// TestPerformanceStatsStrategy tests performance statistics
func TestPerformanceStatsStrategy(t *testing.T) {
	strategy := &PerformanceStatsStrategy{}

	tests := []struct {
		name            string
		posts           []models.Post
		expectedTop     int
		expectedBottom  int
		expectedCount   int
	}{
		{
			name:           "Empty posts",
			posts:          []models.Post{},
			expectedTop:    0,
			expectedBottom: 0,
			expectedCount:  0,
		},
		{
			name: "Single post",
			posts: []models.Post{
				{Likes: 100},
			},
			expectedTop:    100,
			expectedBottom: 100,
			expectedCount:  1,
		},
		{
			name: "Multiple posts with variance",
			posts: []models.Post{
				{Likes: 50},
				{Likes: 100},
				{Likes: 150},
			},
			expectedTop:    150,
			expectedBottom: 50,
			expectedCount:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strategy.Calculate(tt.posts)
			stats := result.(PerformanceStats)

			if stats.TopPostLikes != tt.expectedTop {
				t.Errorf("Expected top %d, got %d", tt.expectedTop, stats.TopPostLikes)
			}
			if stats.BottomPostLikes != tt.expectedBottom {
				t.Errorf("Expected bottom %d, got %d", tt.expectedBottom, stats.BottomPostLikes)
			}
			if stats.PostCount != tt.expectedCount {
				t.Errorf("Expected count %d, got %d", tt.expectedCount, stats.PostCount)
			}
		})
	}
}

// TestStrategyInterface verifies all strategies implement the interface
func TestStrategyInterface(t *testing.T) {
	var aggregateStrat StatisticsStrategy = &AggregateStatsStrategy{}
	var engagementStrat StatisticsStrategy = &EngagementRateStrategy{}
	var performanceStrat StatisticsStrategy = &PerformanceStatsStrategy{}

	posts := []models.Post{{Likes: 100, Comments: 10}}

	if aggregateStrat.Calculate(posts) == nil {
		t.Error("AggregateStatsStrategy.Calculate returned nil")
	}
	if engagementStrat.Calculate(posts) == nil {
		t.Error("EngagementRateStrategy.Calculate returned nil")
	}
	if performanceStrat.Calculate(posts) == nil {
		t.Error("PerformanceStatsStrategy.Calculate returned nil")
	}
}
