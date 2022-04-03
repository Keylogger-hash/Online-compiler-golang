package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"compiler.com/handlers"
)

type Manager struct {
}

func WriteCodeFile(text []byte) {
	currDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile(currDir+"/tmp/main.go", text, fs.FileMode(os.O_WRONLY))
	if err != nil {
		fmt.Println(err)
	}
}
func ExecCode() []byte{
	currDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	text := "package main\nimport \"fmt\"\nfunc main(){\n	fmt.Println(\"I am was writing\")\n}"
	WriteCodeFile([]byte(text))
	pathCompile := currDir + "/tmp/main.go"
	cmd := exec.Command("/opt/homebrew/bin/go", "run", pathCompile)
	output, err := cmd.CombinedOutput()
	return output
}
// Events: [{Message: "Hello, 世界↵1↵", Kind: "stdout", Delay: 0}]
// 0: {Message: "Hello, 世界↵1↵", Kind: "stdout", Delay: 0}
// Delay: 0
// Kind: "stdout"
// Message: "Hello, 世界\n1\n"
func main() {
	mux := http.ServeMux{}
	fs := http.FileServer(http.Dir("static/"))
	mux.Handle("/static/",http.StripPrefix("/static/",fs))
	mux.HandleFunc("/",handlers.HandleIndex)
	mux.HandleFunc("/fmt",handlers.HandleFmt)
	http.ListenAndServe("localhost:8080",&mux)

}
