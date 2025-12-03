package main

import (
	"log"
	"net/http"

	"social-media-analyzer/internal/config"
	"social-media-analyzer/internal/db"
	"social-media-analyzer/internal/db/migrations"
	"social-media-analyzer/internal/service"
	"social-media-analyzer/internal/transport/http/controller"
	"social-media-analyzer/internal/transport/http/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("congif load failed: %v", err)
	}

	db, err := db.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("auto-migrate failed: %v", err)
	}

	// Initialize router
	r := router.New()

	// Initialize services
	vkService := service.NewVKService(&cfg.VK)
	analyticsService := service.NewAnalyticsService(db)
	templateDataService := service.NewTemplateDataService(analyticsService)

	// Initialize controllers
	pageCtrl := controller.NewMainController(templateDataService)
	groupCtrl := controller.NewGroupController(db, vkService)

	// Register routes
	r.GET("/", pageCtrl.GetMainPage)
	r.POST("/api/groups", groupCtrl.AddGroup)

	// Create multiplexer
	mux := http.NewServeMux()
	
	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Use router as fallback for all other routes
	mux.Handle("/", r)
	
	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, mux))
}
