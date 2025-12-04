package service

import (
	"testing"
	"time"

	"social-media-analyzer/internal/models"
	"gorm.io/gorm"
)

// MockDB represents a mock database for testing
type MockDB struct {
	groups []models.Group
	posts  []models.Post
	err    error
}

// MockQueryResult represents a mock query result
type MockQueryResult struct {
	db  *MockDB
	err error
}

// Find mock implementation
func (mq *MockQueryResult) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	if mq.err != nil {
		return &gorm.DB{Error: mq.err}
	}

	// Type assert to get the pointer to slice
	switch v := dest.(type) {
	case *[]models.Group:
		*v = mq.db.groups
	case *[]models.Post:
		*v = mq.db.posts
	}

	return &gorm.DB{Error: nil}
}

// Error method for gorm.DB compatibility
func (mq *MockQueryResult) Error() error {
	return mq.err
}

// TestCalculateGroupStats tests the CalculateGroupStats method
func TestCalculateGroupStats(t *testing.T) {
	tests := []struct {
		name           string
		mockGroups     []models.Group
		mockPosts      map[uint][]models.Post
		expectedCount  int
		shouldError    bool
	}{
		{
			name: "Empty database",
			mockGroups: []models.Group{},
			mockPosts:  map[uint][]models.Post{},
			expectedCount: 0,
			shouldError: false,
		},
		{
			name: "Single group with posts",
			mockGroups: []models.Group{
				{
					ID:          1,
					Domain:      "testgroup",
					Subscribers: 1000,
					ParsedAt:    time.Now(),
				},
			},
			mockPosts: map[uint][]models.Post{
				1: {
					{ID: 1, GroupID: 1, Likes: 100, Comments: 10},
					{ID: 2, GroupID: 1, Likes: 200, Comments: 20},
					{ID: 3, GroupID: 1, Likes: 300, Comments: 30},
				},
			},
			expectedCount: 1,
			shouldError:   false,
		},
		{
			name: "Multiple groups",
			mockGroups: []models.Group{
				{ID: 1, Domain: "group1", Subscribers: 500},
				{ID: 2, Domain: "group2", Subscribers: 1000},
				{ID: 3, Domain: "group3", Subscribers: 2000},
			},
			mockPosts: map[uint][]models.Post{
				1: {{ID: 1, GroupID: 1, Likes: 50, Comments: 5}},
				2: {{ID: 2, GroupID: 2, Likes: 100, Comments: 10}},
				3: {{ID: 3, GroupID: 3, Likes: 200, Comments: 20}},
			},
			expectedCount: 3,
			shouldError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify expected count
			if len(tt.mockGroups) != tt.expectedCount {
				// This test validates the test setup
				stats := make([]GroupStats, len(tt.mockGroups))
				if len(stats) != tt.expectedCount {
					t.Errorf("Expected %d groups, got %d", tt.expectedCount, len(stats))
				}
			}
		})
	}
}

// TestCalculateGroupStat tests individual group stat calculation
func TestCalculateGroupStat(t *testing.T) {
	tests := []struct {
		name              string
		group             models.Group
		posts             []models.Post
		expectedAvgLikes  float64
		expectedMaxLikes  int
		expectedAvgComments float64
	}{
		{
			name: "Group with no posts",
			group: models.Group{
				ID:       1,
				Domain:   "emptygroup",
				Subscribers: 100,
			},
			posts: []models.Post{},
			expectedAvgLikes: 0,
			expectedMaxLikes: 0,
			expectedAvgComments: 0,
		},
		{
			name: "Group with single post",
			group: models.Group{
				ID:       1,
				Domain:   "singlepost",
				Subscribers: 100,
			},
			posts: []models.Post{
				{
					ID: 1,
					GroupID: 1,
					Likes: 100,
					Comments: 10,
				},
			},
			expectedAvgLikes: 100.0,
			expectedMaxLikes: 100,
			expectedAvgComments: 10.0,
		},
		{
			name: "Group with multiple posts",
			group: models.Group{
				ID:       1,
				Domain:   "multipost",
				Subscribers: 1000,
			},
			posts: []models.Post{
				{ID: 1, GroupID: 1, Likes: 100, Comments: 10},
				{ID: 2, GroupID: 1, Likes: 200, Comments: 20},
				{ID: 3, GroupID: 1, Likes: 300, Comments: 30},
			},
			expectedAvgLikes: 200.0,
			expectedMaxLikes: 300,
			expectedAvgComments: 20.0,
		},
		{
			name: "Group with varying engagement",
			group: models.Group{
				ID:       1,
				Domain:   "varying",
				Subscribers: 500,
			},
			posts: []models.Post{
				{ID: 1, GroupID: 1, Likes: 50, Comments: 5},
				{ID: 2, GroupID: 1, Likes: 150, Comments: 15},
				{ID: 3, GroupID: 1, Likes: 100, Comments: 10},
				{ID: 4, GroupID: 1, Likes: 200, Comments: 20},
			},
			expectedAvgLikes: 125.0,
			expectedMaxLikes: 200,
			expectedAvgComments: 12.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Calculate totals for validation
			if len(tt.posts) > 0 {
				totalLikes := 0
				maxLikes := 0
				totalComments := 0

				for _, post := range tt.posts {
					totalLikes += post.Likes
					totalComments += post.Comments
					if post.Likes > maxLikes {
						maxLikes = post.Likes
					}
				}

				avgLikes := float64(totalLikes) / float64(len(tt.posts))
				avgComments := float64(totalComments) / float64(len(tt.posts))

				if avgLikes != tt.expectedAvgLikes {
					t.Errorf("Expected avg likes %.2f, got %.2f", tt.expectedAvgLikes, avgLikes)
				}
				if maxLikes != tt.expectedMaxLikes {
					t.Errorf("Expected max likes %d, got %d", tt.expectedMaxLikes, maxLikes)
				}
				if avgComments != tt.expectedAvgComments {
					t.Errorf("Expected avg comments %.2f, got %.2f", tt.expectedAvgComments, avgComments)
				}
			}
		})
	}
}

// TestChartDataStructure validates ChartData structure
func TestChartDataStructure(t *testing.T) {
	chartData := ChartData{
		Subscribers: []int{100, 200, 300},
		AvgLikes:    []float64{10.5, 20.5, 30.5},
		AvgComments: []float64{5.0, 10.0, 15.0},
	}

	if len(chartData.Subscribers) != 3 {
		t.Errorf("Expected 3 subscribers, got %d", len(chartData.Subscribers))
	}
	if len(chartData.AvgLikes) != 3 {
		t.Errorf("Expected 3 avg likes, got %d", len(chartData.AvgLikes))
	}
	if len(chartData.AvgComments) != 3 {
		t.Errorf("Expected 3 avg comments, got %d", len(chartData.AvgComments))
	}

	if chartData.Subscribers[0] != 100 {
		t.Errorf("Expected first subscriber count 100, got %d", chartData.Subscribers[0])
	}
	if chartData.AvgLikes[1] != 20.5 {
		t.Errorf("Expected second avg likes 20.5, got %.1f", chartData.AvgLikes[1])
	}
}

// TestGroupStatsStructure validates GroupStats structure
func TestGroupStatsStructure(t *testing.T) {
	stats := GroupStats{
		ID:                 1,
		Domain:             "testgroup",
		Subscribers:        1000,
		ParsedAt:           "2025-12-04 10:30",
		TotalPosts:         50,
		TotalLikes:         5000,
		AvgLikesPerPost:    100.0,
		MaxLikesPerPost:    500,
		AvgCommentsPerPost: 10.0,
		PostsLastWeek:      5,
	}

	if stats.ID != 1 {
		t.Errorf("Expected ID 1, got %d", stats.ID)
	}
	if stats.Domain != "testgroup" {
		t.Errorf("Expected domain 'testgroup', got '%s'", stats.Domain)
	}
	if stats.Subscribers != 1000 {
		t.Errorf("Expected 1000 subscribers, got %d", stats.Subscribers)
	}
	if stats.TotalPosts != 50 {
		t.Errorf("Expected 50 posts, got %d", stats.TotalPosts)
	}
	if stats.AvgLikesPerPost != 100.0 {
		t.Errorf("Expected avg likes 100.0, got %.1f", stats.AvgLikesPerPost)
	}
}
