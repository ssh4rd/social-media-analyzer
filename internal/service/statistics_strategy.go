package service

import (
	"social-media-analyzer/internal/models"
)

// StatisticsStrategy defines the interface for different statistics calculation strategies
type StatisticsStrategy interface {
	Calculate(posts []models.Post) interface{}
}

// AggregateStatsStrategy calculates aggregate statistics (likes, comments, max, etc.)
type AggregateStatsStrategy struct{}

// AggregateStats represents aggregate statistics for posts
type AggregateStats struct {
	TotalLikes         int
	AvgLikesPerPost    float64
	MaxLikesPerPost    int
	TotalComments      int
	AvgCommentsPerPost float64
}

// Calculate implements StatisticsStrategy for aggregate statistics
func (s *AggregateStatsStrategy) Calculate(posts []models.Post) interface{} {
	if len(posts) == 0 {
		return AggregateStats{
			TotalLikes:         0,
			AvgLikesPerPost:    0,
			MaxLikesPerPost:    0,
			TotalComments:      0,
			AvgCommentsPerPost: 0,
		}
	}

	totalLikes := 0
	maxLikes := 0
	totalComments := 0

	for _, post := range posts {
		totalLikes += post.Likes
		totalComments += post.Comments
		if post.Likes > maxLikes {
			maxLikes = post.Likes
		}
	}

	return AggregateStats{
		TotalLikes:         totalLikes,
		AvgLikesPerPost:    float64(totalLikes) / float64(len(posts)),
		MaxLikesPerPost:    maxLikes,
		TotalComments:      totalComments,
		AvgCommentsPerPost: float64(totalComments) / float64(len(posts)),
	}
}

// EngagementRateStrategy calculates engagement rates and ratios
type EngagementRateStrategy struct{}

// EngagementRate represents engagement metrics
type EngagementRate struct {
	LikeToCommentRatio float64 // Likes per comment
	AvgViews           int
	TotalEngagement    int
}

// Calculate implements StatisticsStrategy for engagement rates
func (s *EngagementRateStrategy) Calculate(posts []models.Post) interface{} {
	if len(posts) == 0 {
		return EngagementRate{
			LikeToCommentRatio: 0,
			AvgViews:           0,
			TotalEngagement:    0,
		}
	}

	totalLikes := 0
	totalComments := 0
	totalViews := 0
	totalEngagement := 0

	for _, post := range posts {
		totalLikes += post.Likes
		totalComments += post.Comments
		totalViews += post.Views
		totalEngagement += post.Likes + post.Comments + post.Reactions
	}

	likeToCommentRatio := 0.0
	if totalComments > 0 {
		likeToCommentRatio = float64(totalLikes) / float64(totalComments)
	}

	return EngagementRate{
		LikeToCommentRatio: likeToCommentRatio,
		AvgViews:           totalViews / len(posts),
		TotalEngagement:    totalEngagement,
	}
}

// PerformanceStatsStrategy calculates performance-based statistics
type PerformanceStatsStrategy struct{}

// PerformanceStats represents performance metrics
type PerformanceStats struct {
	TopPostLikes    int
	BottomPostLikes int
	VarianceInLikes float64 // Measure of consistency
	PostCount       int
}

// Calculate implements StatisticsStrategy for performance statistics
func (s *PerformanceStatsStrategy) Calculate(posts []models.Post) interface{} {
	if len(posts) == 0 {
		return PerformanceStats{
			TopPostLikes:    0,
			BottomPostLikes: 0,
			VarianceInLikes: 0,
			PostCount:       0,
		}
	}

	topLikes := 0
	bottomLikes := posts[0].Likes
	sumSquares := 0
	avgLikes := 0

	// Calculate totals
	totalLikes := 0
	for _, post := range posts {
		totalLikes += post.Likes
		if post.Likes > topLikes {
			topLikes = post.Likes
		}
		if post.Likes < bottomLikes {
			bottomLikes = post.Likes
		}
	}

	avgLikes = totalLikes / len(posts)

	// Calculate variance
	for _, post := range posts {
		diff := post.Likes - avgLikes
		sumSquares += diff * diff
	}

	variance := float64(sumSquares) / float64(len(posts))

	return PerformanceStats{
		TopPostLikes:    topLikes,
		BottomPostLikes: bottomLikes,
		VarianceInLikes: variance,
		PostCount:       len(posts),
	}
}
