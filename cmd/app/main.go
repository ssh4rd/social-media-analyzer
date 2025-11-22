package main

import (
	"log"
	"net/http"

	"social-media-analyzer/internal/config"
	"social-media-analyzer/internal/db"
	"social-media-analyzer/internal/db/migrations"
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

	// Initialize controller
	pageCtrl := &controller.MainController{}
	employeeCtrl := &controller.EmployeeController{}

	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Register routes
	r.GET("/", pageCtrl.GetMainPage)
	r.GET("/employee/:id", employeeCtrl.GetEmployee)

	// Start server with both router and static handler
	http.Handle("/", r)
	http.ListenAndServe(":"+cfg.Server.Port, nil)
}
