package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

const Root string = "C:\\Users\\1\\go\\src\\Online-compiler-golang\\public\\templates"
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	filenames := []string{Root+"\\index.html"}
	tpl, err := template.ParseFiles(filenames...)
	if err != nil {
		fmt.Println(err)
	}
	tpl.Execute(w, nil)
}
