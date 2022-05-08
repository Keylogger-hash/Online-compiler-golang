package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"compiler.com/handlers"
)

// Events: [{Message: "Hello, 世界↵1↵", Kind: "stdout", Delay: 0}]
// 0: {Message: "Hello, 世界↵1↵", Kind: "stdout", Delay: 0}
// Delay: 0
// Kind: "stdout"
// Message: "Hello, 世界\n1\n"
func main() {
	mux := mux.NewRouter()
	fs := http.FileServer(http.Dir("./static/"))
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/",fs))
	// mux.PathPrefix("/static/").Handler(fs)
	//mux.Path("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", handlers.HandleIndex)
	mux.HandleFunc("/compile", handlers.HandleCompile)
	mux.HandleFunc("/contacts",handlers.HandleContacts)
	mux.HandleFunc("/about",handlers.HandleAbout)
	mux.HandleFunc("/terms",handlers.HandleTerms)
	mux.HandleFunc("/share",handlers.HandleShare)
	mux.HandleFunc("/share/p/{id}",handlers.HandleShareInfo)
	mux.HandleFunc("/why",handlers.HandleWhy)
	mux.HandleFunc("/fmt", handlers.HandleFmt)
	fmt.Println("Starting server...")
	fmt.Println("Listen and serve on the port 8080")
	http.ListenAndServe("localhost:8080", mux)

}
