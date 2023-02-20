package handlers

import (
	"finances_manager_go/utils"
	"html/template"
	"net/http"

	"gorm.io/gorm"
)

type APIEnv struct {
	DB *gorm.DB
}

func (a *APIEnv) AddIncome(w http.ResponseWriter, r *http.Request) {
	utils.MethodValidation(w, r, "POST")
}

var tpl = template.Must(template.ParseFiles("index.html"))

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}
