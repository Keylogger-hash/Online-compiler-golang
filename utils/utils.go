package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/tools/imports"
)

type BuildResults struct {
	pathCompile string
	pathOutput  string
	Data        []byte
	Stdout      io.ReadCloser
	Stderr      io.ReadCloser
	Errors      error
	Events      []string
}

func NewBuildResult() *BuildResults {
	return &BuildResults{}
}
func CheckCodePackageIsMain(build *BuildResults) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", build.Data, parser.AllErrors)
	if err != nil {
		build.Errors = err
		return
	}
	if file.Name.Name != "main" {
		build.Errors = fmt.Errorf("package is not main")
		return
	}
	return
}
func FormatFmt(build *BuildResults, body string) {
	dest, err := format.Source([]byte(body))
	if err != nil {
		build.Errors = err
		return
	}
	finishImports, err := imports.Process("", dest, nil)
	if err != nil {
		build.Errors = err
		return
	}
	build.Data = finishImports
	return
}
func WriteCodeFile(build *BuildResults) {
	uid := uuid.New().String()
	nameContainer := "main_" + uid
	build.pathOutput = fmt.Sprintf("/tmp/build/%s/", uid)
	build.pathCompile = build.pathOutput + nameContainer
	err := os.Mkdir(build.pathOutput, 0750)
	if err != nil {
		build.Errors = err
		return
	}
	file, err := os.Create(build.pathCompile + ".go")
	if err != nil {
		build.Errors = err
		return
	}
	_, err = file.Write(build.Data)
	if err != nil {
		build.Errors = err
	}
	return
}
func CompileCode(build *BuildResults) {
	cmd := exec.Command("go", "build", "-o", build.pathOutput, build.pathCompile+".go")
	fmt.Println(build.pathOutput, build.pathCompile+".go")
	var stderr strings.Builder
	cmd.Env = append(cmd.Env, "GOOS=linux")
	cmd.Env = append(cmd.Env, "GOARCH=arm64")
	cmd.Env = append(cmd.Env, "GOPATH=/Users/pavelmorozov/go")
	cmd.Env = append(cmd.Env, "GOCACHE=/Users/pavelmorozov/Library/Caches/go-build")
	cmd.Stderr = &stderr
	err := cmd.Run()
	fmt.Println(stderr.String())
	if err != nil {
		fmt.Println(stderr.String())
		build.Errors = fmt.Errorf("%s", stderr.String())
	} else {
		build.Errors = nil
	}
}
func EncodeBinaryFile(build *BuildResults) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	defer encoder.Close()
	b, err := ioutil.ReadFile(build.pathCompile)
	if err != nil {
		return nil, err
	}
	encoder.Write(b)
	return &buf, nil
}

// func main() {
// 	const data = `package main
// 	func main(){
// 		fmt.Println(10)
// 	}
// 	`
// 	f, err := FormatFmt(data)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	ok, err := CheckCodePackageIsMain(f)
// 	if ok {
// 		fmt.Println("package is main")
// 	} else {
// 		fmt.Println("package not main")
// 	}
// 	p, o, err := WriteCodeFile(f)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	CompileCode(p, o)
// 	buf, err := EncodeBinaryFile(p)
// 	dst, err := DecodeBinaryFile(buf)
// 	CreateExecutableFile(dst)
// 	// CompileCode()

// }
