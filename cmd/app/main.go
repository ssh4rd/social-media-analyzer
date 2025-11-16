package main

import (
	"log"
	"net/http"
	"os"

	"TP_Andreev/internal/db"
	"TP_Andreev/internal/db/migrations"
	"TP_Andreev/internal/transport/http/controller"
	"TP_Andreev/internal/transport/http/router"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	db, err := db.Connect()
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

	// Register routes
	r.GET("/", pageCtrl.GetMainPage)

	// Start server
	http.ListenAndServe(":"+port, r)
}
