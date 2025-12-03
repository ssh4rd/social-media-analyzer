package service

import (
	"log"

	"social-media-analyzer/internal/models"
	"gorm.io/gorm"
)

type GroupStats struct {
	ID                 uint
	Domain             string
	Subscribers        int
	ParsedAt           string
	TotalPosts         int
	TotalLikes         int
	AvgLikesPerPost    float64
	MaxLikesPerPost    int
	AvgCommentsPerPost float64
	PostsLastWeek      int
}

type ChartData struct {
	Subscribers []int     `json:"subscribers"`
	AvgLikes    []float64 `json:"avgLikes"`
	AvgComments []float64 `json:"avgComments"`
}

type AnalyticsService struct {
	db *gorm.DB
}

func NewAnalyticsService(db *gorm.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

// CalculateGroupStats calculates statistics for all groups
func (as *AnalyticsService) CalculateGroupStats() ([]GroupStats, error) {
	var groups []models.Group
	if err := as.db.Find(&groups).Error; err != nil {
		return nil, err
	}

	var stats []GroupStats
	for _, group := range groups {
		groupStat := as.calculateGroupStat(group)
		stats = append(stats, groupStat)
	}

	return stats, nil
}

// calculateGroupStat calculates statistics for a single group
func (as *AnalyticsService) calculateGroupStat(group models.Group) GroupStats {
	var posts []models.Post
	if err := as.db.Where("group_id = ?", group.ID).Find(&posts).Error; err != nil {
		log.Printf("Failed to fetch posts for group %s: %v\n", group.Domain, err)
	}

	parsedAt := "-"
	if !group.ParsedAt.IsZero() {
		parsedAt = group.ParsedAt.Format("2006-01-02 15:04")
	}

	stats := GroupStats{
		ID:         group.ID,
		Domain:     group.Domain,
		Subscribers: group.Subscribers,
		ParsedAt:   parsedAt,
		TotalPosts: len(posts),
	}

	if len(posts) > 0 {
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

		stats.TotalLikes = totalLikes
		stats.MaxLikesPerPost = maxLikes
		stats.AvgLikesPerPost = float64(totalLikes) / float64(len(posts))
		stats.AvgCommentsPerPost = float64(totalComments) / float64(len(posts))
	}

	return stats
}

// CalculateChartData calculates data for charts (dependence of likes/comments on subscribers)
func (as *AnalyticsService) CalculateChartData() (ChartData, error) {
	stats, err := as.CalculateGroupStats()
	if err != nil {
		return ChartData{}, err
	}

	chartData := ChartData{
		Subscribers: make([]int, len(stats)),
		AvgLikes:    make([]float64, len(stats)),
		AvgComments: make([]float64, len(stats)),
	}

	for i, stat := range stats {
		chartData.Subscribers[i] = stat.Subscribers
		chartData.AvgLikes[i] = stat.AvgLikesPerPost
		chartData.AvgComments[i] = stat.AvgCommentsPerPost
	}

	return chartData, nil
}
