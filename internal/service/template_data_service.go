package service

// TemplateDataService handles data preparation for template rendering
type TemplateDataService struct {
	analyticsService *AnalyticsService
}

// TemplateGroupData represents formatted group data for template rendering
type TemplateGroupData struct {
	ID                 uint
	Domain             string
	Subscribers        int
	ParsedAt           string
	TotalPosts         int
	TotalLikes         int
	AvgLikesPerPost    float64
	MaxLikesPerPost    int
	AvgCommentsPerPost float64
}

// ChartDataForTemplate represents chart data formatted for template rendering
type ChartDataForTemplate struct {
	Subscribers []int     `json:"subscribers"`
	AvgLikes    []float64 `json:"avgLikes"`
	AvgComments []float64 `json:"avgComments"`
}

func NewTemplateDataService(analyticsService *AnalyticsService) *TemplateDataService {
	return &TemplateDataService{analyticsService: analyticsService}
}

// PrepareGroupsForTemplate prepares group statistics for template rendering
func (tds *TemplateDataService) PrepareGroupsForTemplate() ([]TemplateGroupData, error) {
	stats, err := tds.analyticsService.CalculateGroupStats()
	if err != nil {
		return nil, err
	}

	templateData := make([]TemplateGroupData, len(stats))
	for i, stat := range stats {
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

	return templateData, nil
}

// PrepareChartDataForTemplate prepares chart data for template rendering
func (tds *TemplateDataService) PrepareChartDataForTemplate() (ChartDataForTemplate, error) {
	chartData, err := tds.analyticsService.CalculateChartData()
	if err != nil {
		return ChartDataForTemplate{}, err
	}

	return ChartDataForTemplate{
		Subscribers: chartData.Subscribers,
		AvgLikes:    chartData.AvgLikes,
		AvgComments: chartData.AvgComments,
	}, nil
}
