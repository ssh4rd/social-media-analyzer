package controller

import (
	"html/template"
	"net/http"

	"TP_Andreev/internal/transport/http/router"
)

var employeeTpl = template.Must(template.ParseFiles("web/templates/employee.html"))

type EmployeeController struct{}

// GetEmployee handles GET /employees/:id requests
func (ec *EmployeeController) GetEmployee(w http.ResponseWriter, r *http.Request, params router.Params) {
	// id := params["id"]
	// TODO: Fetch employee by ID
	employeeTpl.Execute(w, nil)
}
