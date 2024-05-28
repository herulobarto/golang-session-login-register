package controllers

import (
	"fmt"
	"html/template"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {

	temp, _ := template.ParseFiles("views/Index.html")
	temp.Execute(w, nil)

}

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		temp, _ := template.ParseFiles("views/Login.html")
		temp.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		// proses login
		fmt.Println("Proses Login")
	}

}
