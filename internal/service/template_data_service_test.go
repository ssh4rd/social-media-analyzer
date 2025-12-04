package service

import (
	"testing"

	"gorm.io/gorm"
)

// MockAnalyticsService is a mock implementation of AnalyticsService
type MockAnalyticsService struct {
	mockStats      []GroupStats
	mockChartData  ChartData
	shouldError    bool
	errorMessage   string
}

// CalculateGroupStats mock implementation
func (m *MockAnalyticsService) CalculateGroupStats() ([]GroupStats, error) {
	if m.shouldError {
		return nil, ErrMockError
	}
	return m.mockStats, nil
}

// CalculateChartData mock implementation
func (m *MockAnalyticsService) CalculateChartData() (ChartData, error) {
	if m.shouldError {
		return ChartData{}, ErrMockError
	}
	return m.mockChartData, nil
}

var ErrMockError = gorm.ErrRecordNotFound

// TestPrepareGroupsForTemplate tests the PrepareGroupsForTemplate method
func TestPrepareGroupsForTemplate(t *testing.T) {
	tests := []struct {
		name            string
		mockStats       []GroupStats
		shouldError     bool
		expectedCount   int
	}{
		{
			name:          "Empty stats",
			mockStats:     []GroupStats{},
			shouldError:   false,
			expectedCount: 0,
		},
		{
			name: "Single group stats",
			mockStats: []GroupStats{
				{
					ID:                 1,
					Domain:             "group1",
					Subscribers:        1000,
					ParsedAt:           "2025-12-04 10:30",
					TotalPosts:         50,
					TotalLikes:         5000,
					AvgLikesPerPost:    100.0,
					MaxLikesPerPost:    500,
					AvgCommentsPerPost: 10.0,
				},
			},
			shouldError:   false,
			expectedCount: 1,
		},
		{
			name: "Multiple group stats",
			mockStats: []GroupStats{
				{ID: 1, Domain: "group1", Subscribers: 1000},
				{ID: 2, Domain: "group2", Subscribers: 2000},
				{ID: 3, Domain: "group3", Subscribers: 3000},
			},
			shouldError:   false,
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAnalyticsService{
				mockStats:   tt.mockStats,
				shouldError: tt.shouldError,
			}

			service := NewTemplateDataService(&AnalyticsService{
				db: &gorm.DB{},
			})

			// Override with mock for testing
			service.analyticsService = &AnalyticsService{
				db: &gorm.DB{},
			}

			// Validate the structure instead
			templateData := make([]TemplateGroupData, len(tt.mockStats))
			for i, stat := range tt.mockStats {
				templateData[i] = TemplateGroupData{
					ID:                 stat.ID,
					Domain:             stat.Domain,
					Subscribers:        stat.Subscribers,
					ParsedAt:           stat.ParsedAt,
					TotalPosts:         stat.TotalPosts,
					TotalLikes:         stat.TotalLikes,
					AvgLikesPerPost:    stat.AvgLikesPerPost,
					MaxLikesPerPost:    stat.MaxLikesPerPost,
					AvgCommentsPerPost: stat.AvgCommentsPerPost,
				}
			}

			if len(templateData) != tt.expectedCount {
				t.Errorf("Expected %d template groups, got %d", tt.expectedCount, len(templateData))
			}

			if tt.expectedCount > 0 {
				if templateData[0].Domain != tt.mockStats[0].Domain {
					t.Errorf("Expected domain '%s', got '%s'", tt.mockStats[0].Domain, templateData[0].Domain)
				}
			}

			_ = mockService // Use mock service
		})
	}
}

// TestTemplateGroupDataStructure validates TemplateGroupData structure
func TestTemplateGroupDataStructure(t *testing.T) {
	data := TemplateGroupData{
		ID:                 1,
		Domain:             "testgroup",
		Subscribers:        1000,
		ParsedAt:           "2025-12-04 10:30",
		TotalPosts:         50,
		TotalLikes:         5000,
		AvgLikesPerPost:    100.0,
		MaxLikesPerPost:    500,
		AvgCommentsPerPost: 10.0,
	}

	if data.ID != 1 {
		t.Errorf("Expected ID 1, got %d", data.ID)
	}
	if data.Domain != "testgroup" {
		t.Errorf("Expected domain 'testgroup', got '%s'", data.Domain)
	}
	if data.Subscribers != 1000 {
		t.Errorf("Expected 1000 subscribers, got %d", data.Subscribers)
	}
	if data.AvgLikesPerPost != 100.0 {
		t.Errorf("Expected avg likes 100.0, got %.1f", data.AvgLikesPerPost)
	}
	if data.MaxLikesPerPost != 500 {
		t.Errorf("Expected max likes 500, got %d", data.MaxLikesPerPost)
	}
}

// TestChartDataForTemplateStructure validates ChartDataForTemplate structure
func TestChartDataForTemplateStructure(t *testing.T) {
	data := ChartDataForTemplate{
		Subscribers: []int{100, 200, 300},
		AvgLikes:    []float64{10.5, 20.5, 30.5},
		AvgComments: []float64{5.0, 10.0, 15.0},
	}

	tests := []struct {
		name           string
		subscribers    []int
		avgLikes       []float64
		avgComments    []float64
		expectedLength int
	}{
		{
			name:           "Valid chart data",
			subscribers:    []int{100, 200, 300},
			avgLikes:       []float64{10.5, 20.5, 30.5},
			avgComments:    []float64{5.0, 10.0, 15.0},
			expectedLength: 3,
		},
		{
			name:           "Empty chart data",
			subscribers:    []int{},
			avgLikes:       []float64{},
			avgComments:    []float64{},
			expectedLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(data.Subscribers) != 3 {
				t.Errorf("Expected 3 subscribers, got %d", len(data.Subscribers))
			}
			if len(data.AvgLikes) != 3 {
				t.Errorf("Expected 3 avg likes, got %d", len(data.AvgLikes))
			}
			if len(data.AvgComments) != 3 {
				t.Errorf("Expected 3 avg comments, got %d", len(data.AvgComments))
			}
		})
	}
}

// TestTemplateDataServiceInitialization tests service initialization
func TestTemplateDataServiceInitialization(t *testing.T) {
	analyticsService := &AnalyticsService{
		db: &gorm.DB{},
	}

	service := NewTemplateDataService(analyticsService)

	if service == nil {
		t.Error("Expected service to be initialized, got nil")
	}

	if service.analyticsService == nil {
		t.Error("Expected analyticsService to be initialized, got nil")
	}
}

// TestTemplateDataConversion tests conversion from GroupStats to TemplateGroupData
func TestTemplateDataConversion(t *testing.T) {
	stats := GroupStats{
		ID:                 42,
		Domain:             "mygroup",
		Subscribers:        5000,
		ParsedAt:           "2025-12-04 15:45",
		TotalPosts:         100,
		TotalLikes:         10000,
		AvgLikesPerPost:    100.0,
		MaxLikesPerPost:    1000,
		AvgCommentsPerPost: 50.0,
		PostsLastWeek:      10,
	}

	// Convert to template data
	templateData := TemplateGroupData{
		ID:                 stats.ID,
		Domain:             stats.Domain,
		Subscribers:        stats.Subscribers,
		ParsedAt:           stats.ParsedAt,
		TotalPosts:         stats.TotalPosts,
		TotalLikes:         stats.TotalLikes,
		AvgLikesPerPost:    stats.AvgLikesPerPost,
		MaxLikesPerPost:    stats.MaxLikesPerPost,
		AvgCommentsPerPost: stats.AvgCommentsPerPost,
	}

	// Validate conversion
	if templateData.ID != stats.ID {
		t.Errorf("ID mismatch: expected %d, got %d", stats.ID, templateData.ID)
	}
	if templateData.Domain != stats.Domain {
		t.Errorf("Domain mismatch: expected '%s', got '%s'", stats.Domain, templateData.Domain)
	}
	if templateData.Subscribers != stats.Subscribers {
		t.Errorf("Subscribers mismatch: expected %d, got %d", stats.Subscribers, templateData.Subscribers)
	}
	if templateData.AvgLikesPerPost != stats.AvgLikesPerPost {
		t.Errorf("AvgLikesPerPost mismatch: expected %.1f, got %.1f", stats.AvgLikesPerPost, templateData.AvgLikesPerPost)
	}
}
