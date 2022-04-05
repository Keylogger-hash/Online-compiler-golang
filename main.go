package main

import (
	"net/http"

	"compiler.com/handlers"
)

// Events: [{Message: "Hello, 世界↵1↵", Kind: "stdout", Delay: 0}]
// 0: {Message: "Hello, 世界↵1↵", Kind: "stdout", Delay: 0}
// Delay: 0
// Kind: "stdout"
// Message: "Hello, 世界\n1\n"
func main() {
	mux := http.ServeMux{}
	fs := http.FileServer(http.Dir("static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", handlers.HandleIndex)
	mux.HandleFunc("/compile", handlers.HandleCompile)
	mux.HandleFunc("/fmt", handlers.HandleFmt)
	http.ListenAndServe("localhost:8080", &mux)

}
