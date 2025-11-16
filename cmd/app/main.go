package main

import (
	"net/http"
	"os"

	"TP_Andreev/internal/transport/http/controller"
	"TP_Andreev/internal/transport/http/router"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
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
	http.ListenAndServe(":"+port, nil)
}
