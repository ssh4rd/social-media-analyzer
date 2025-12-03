package controller

import (
	"encoding/json"
	"html/template"
	"net/http"

	"social-media-analyzer/internal/service"
	"social-media-analyzer/internal/transport/http/router"
)

type MainController struct {
	templateDataService *service.TemplateDataService
}

type PageData struct {
	Groups     []service.TemplateGroupData
	ChartData  service.ChartDataForTemplate
}

func NewMainController(templateDataService *service.TemplateDataService) *MainController {
	return &MainController{templateDataService: templateDataService}
}

// GetMainPage handles GET / requests
func (mc *MainController) GetMainPage(w http.ResponseWriter, r *http.Request, params router.Params) {
	// Prepare group data for template rendering
	groupData, err := mc.templateDataService.PrepareGroupsForTemplate()
	if err != nil {
		http.Error(w, "Failed to prepare template data", http.StatusInternalServerError)
		return
	}

	// Prepare chart data for template rendering
	chartData, err := mc.templateDataService.PrepareChartDataForTemplate()
	if err != nil {
		http.Error(w, "Failed to prepare chart data", http.StatusInternalServerError)
		return
	}

	pageData := PageData{
		Groups:    groupData,
		ChartData: chartData,
	}

	funcMap := template.FuncMap{
		"json": func(v interface{}) template.JS {
			b, _ := json.Marshal(v)
			return template.JS(b)
		},
	}

	tpl := template.Must(template.New("main.html").Funcs(funcMap).ParseFiles("web/templates/main.html"))
	tpl.Execute(w, pageData)
}
