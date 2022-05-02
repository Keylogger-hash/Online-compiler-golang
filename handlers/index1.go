package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

//const Root string = "/Users/pavelmorozov/go/src/golang-online-compiler/public/templates"
func HandleIndex1(w http.ResponseWriter, r *http.Request) {
	filenames := []string{Root+"/index1.html"}
	tpl, err := template.ParseFiles(filenames...)
	if err != nil {
		fmt.Println(err)
	}
	tpl.Execute(w, nil)
}