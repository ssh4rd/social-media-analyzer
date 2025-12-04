package service

import (
	"social-media-analyzer/internal/config"
	"gorm.io/gorm"
)

// ServiceFactory creates and configures service instances
// This implements the Factory design pattern
type ServiceFactory struct {
	config *config.Config
	db     *gorm.DB
}

// ServiceContainer holds all initialized services
type ServiceContainer struct {
	VKService             *VKService
	AnalyticsService      *AnalyticsService
	TemplateDataService   *TemplateDataService
	AggregateStrategy     StatisticsStrategy
	EngagementStrategy    StatisticsStrategy
	PerformanceStrategy   StatisticsStrategy
}

// NewServiceFactory creates a new service factory
func NewServiceFactory(cfg *config.Config, db *gorm.DB) *ServiceFactory {
	return &ServiceFactory{
		config: cfg,
		db:     db,
	}
}

// CreateServices creates and initializes all application services
func (sf *ServiceFactory) CreateServices() *ServiceContainer {
	// Create core services
	vkService := sf.createVKService()
	analyticsService := sf.createAnalyticsService()
	templateDataService := sf.createTemplateDataService(analyticsService)

	// Create statistics strategies
	aggregateStrategy := &AggregateStatsStrategy{}
	engagementStrategy := &EngagementRateStrategy{}
	performanceStrategy := &PerformanceStatsStrategy{}

	return &ServiceContainer{
		VKService:           vkService,
		AnalyticsService:    analyticsService,
		TemplateDataService: templateDataService,
		AggregateStrategy:   aggregateStrategy,
		EngagementStrategy:  engagementStrategy,
		PerformanceStrategy: performanceStrategy,
	}
}

// createVKService creates and configures VK API service
func (sf *ServiceFactory) createVKService() *VKService {
	return NewVKService(&sf.config.VK)
}

// createAnalyticsService creates and configures analytics service
func (sf *ServiceFactory) createAnalyticsService() *AnalyticsService {
	return NewAnalyticsService(sf.db)
}

// createTemplateDataService creates and configures template data service
func (sf *ServiceFactory) createTemplateDataService(analyticsService *AnalyticsService) *TemplateDataService {
	return NewTemplateDataService(analyticsService)
}

// CreateVKServiceOnly creates only the VK service (useful for testing or specific use cases)
func (sf *ServiceFactory) CreateVKServiceOnly() *VKService {
	return sf.createVKService()
}

// CreateAnalyticsServicesOnly creates analytics service with strategies
func (sf *ServiceFactory) CreateAnalyticsServicesOnly() (*AnalyticsService, StatisticsStrategy, StatisticsStrategy, StatisticsStrategy) {
	analyticsService := sf.createAnalyticsService()
	return analyticsService, &AggregateStatsStrategy{}, &EngagementRateStrategy{}, &PerformanceStatsStrategy{}
}
