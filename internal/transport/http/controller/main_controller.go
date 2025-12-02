package controller

import (
	"html/template"
	"net/http"

	"social-media-analyzer/internal/transport/http/router"
)

type MainController struct{}

// GetMainPage handles GET / requests
func (pc *MainController) GetMainPage(w http.ResponseWriter, r *http.Request, params router.Params) {
	tpl := template.Must(template.ParseFiles("web/templates/main.html"))
	tpl.Execute(w, nil)
}
