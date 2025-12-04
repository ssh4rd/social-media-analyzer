package service

import (
	"testing"

	"social-media-analyzer/internal/config"
	"gorm.io/gorm"
)

// TestServiceFactoryCreation tests factory initialization
func TestServiceFactoryCreation(t *testing.T) {
	cfg := &config.Config{
		VK: config.VKConfig{
			AccessToken: "test_token",
			APIVersion:  "5.131",
		},
	}
	db := &gorm.DB{}

	factory := NewServiceFactory(cfg, db)

	if factory == nil {
		t.Error("Expected factory to be initialized, got nil")
	}
	if factory.config == nil {
		t.Error("Expected config to be set in factory")
	}
	if factory.db == nil {
		t.Error("Expected db to be set in factory")
	}
}

// TestServiceContainerCreation tests full service container creation
func TestServiceContainerCreation(t *testing.T) {
	cfg := &config.Config{
		VK: config.VKConfig{
			AccessToken: "test_token",
			APIVersion:  "5.131",
		},
	}
	db := &gorm.DB{}

	factory := NewServiceFactory(cfg, db)
	services := factory.CreateServices()

	if services == nil {
		t.Fatal("Expected services container to be initialized, got nil")
	}

	if services.VKService == nil {
		t.Error("Expected VKService to be initialized")
	}
	if services.AnalyticsService == nil {
		t.Error("Expected AnalyticsService to be initialized")
	}
	if services.TemplateDataService == nil {
		t.Error("Expected TemplateDataService to be initialized")
	}
}

// TestServiceContainerStrategies tests strategy initialization
func TestServiceContainerStrategies(t *testing.T) {
	cfg := &config.Config{
		VK: config.VKConfig{
			AccessToken: "test_token",
			APIVersion:  "5.131",
		},
	}
	db := &gorm.DB{}

	factory := NewServiceFactory(cfg, db)
	services := factory.CreateServices()

	if services.AggregateStrategy == nil {
		t.Error("Expected AggregateStrategy to be initialized")
	}
	if services.EngagementStrategy == nil {
		t.Error("Expected EngagementStrategy to be initialized")
	}
	if services.PerformanceStrategy == nil {
		t.Error("Expected PerformanceStrategy to be initialized")
	}

	// Verify strategies are of correct type
	if _, ok := services.AggregateStrategy.(*AggregateStatsStrategy); !ok {
		t.Error("Expected AggregateStrategy to be *AggregateStatsStrategy")
	}
	if _, ok := services.EngagementStrategy.(*EngagementRateStrategy); !ok {
		t.Error("Expected EngagementStrategy to be *EngagementRateStrategy")
	}
	if _, ok := services.PerformanceStrategy.(*PerformanceStatsStrategy); !ok {
		t.Error("Expected PerformanceStrategy to be *PerformanceStatsStrategy")
	}
}

// TestFactoryCreateVKServiceOnly tests creating only VK service
func TestFactoryCreateVKServiceOnly(t *testing.T) {
	cfg := &config.Config{
		VK: config.VKConfig{
			AccessToken: "test_token",
			APIVersion:  "5.131",
		},
	}
	db := &gorm.DB{}

	factory := NewServiceFactory(cfg, db)
	vkService := factory.CreateVKServiceOnly()

	if vkService == nil {
		t.Error("Expected VKService to be created")
	}
}

// TestFactoryCreateAnalyticsServicesOnly tests creating analytics service with strategies
func TestFactoryCreateAnalyticsServicesOnly(t *testing.T) {
	cfg := &config.Config{
		VK: config.VKConfig{
			AccessToken: "test_token",
			APIVersion:  "5.131",
		},
	}
	db := &gorm.DB{}

	factory := NewServiceFactory(cfg, db)
	analyticsService, aggStrat, engStrat, perfStrat := factory.CreateAnalyticsServicesOnly()

	if analyticsService == nil {
		t.Error("Expected AnalyticsService to be created")
	}
	if aggStrat == nil {
		t.Error("Expected AggregateStatsStrategy to be created")
	}
	if engStrat == nil {
		t.Error("Expected EngagementRateStrategy to be created")
	}
	if perfStrat == nil {
		t.Error("Expected PerformanceStatsStrategy to be created")
	}
}

// TestFactoryMultipleCreations tests that factory creates independent instances
func TestFactoryMultipleCreations(t *testing.T) {
	cfg := &config.Config{
		VK: config.VKConfig{
			AccessToken: "test_token",
			APIVersion:  "5.131",
		},
	}
	db := &gorm.DB{}

	factory := NewServiceFactory(cfg, db)

	services1 := factory.CreateServices()
	services2 := factory.CreateServices()

	// Services should be different instances
	if services1.VKService == services2.VKService {
		t.Error("Expected different VKService instances")
	}
	if services1.AnalyticsService == services2.AnalyticsService {
		t.Error("Expected different AnalyticsService instances")
	}
	if services1.TemplateDataService == services2.TemplateDataService {
		t.Error("Expected different TemplateDataService instances")
	}
}
