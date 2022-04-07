package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/google/uuid"
	"golang.org/x/tools/imports"
)

type Code struct {
}

func CheckCodePackageIsMain(fdata []byte) (bool, error) {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "", fdata, parser.AllErrors)
	if err != nil {
		return false, err
	}
	if file.Name.Name != "main" {
		return false, nil
	}
	return true, nil
}
func FormatFmt(body string) ([]byte, error) {
	dest, err := format.Source([]byte(body))
	if err != nil {
		return nil, err
	}
	finishImports, err := imports.Process("", dest, nil)
	if err != nil {
		return nil, err
	}
	return finishImports, nil
}
func WriteCodeFile(data []byte) (string, string, error) {
	uid := uuid.New().String()
	nameContainer := "main_" + uid
	pathName := fmt.Sprintf("tmp/%s/", uid)
	pathCompile := pathName + nameContainer
	err := os.Mkdir(pathName, fs.ModeTemporary)
	if err != nil {
		return "", "", err
	}
	file, err := os.Create(pathCompile + ".go")
	if err != nil {
		return "", "", err
	}
	_, err = file.Write(data)
	if err != nil {
		return "", "", err
	}
	return pathCompile, pathName, nil
}
func CompileCode(pathCompile string, pathOutput string) {
	cmd := exec.Command("go", "build", "-o", pathOutput, pathCompile+".go")
	cmd.Env = append(cmd.Env, "GOOS=linux")
	cmd.Env = append(cmd.Env, "GOARCH=amd64")
	cmd.Env = append(cmd.Env, "GOPATH=C:\\Users\\1\\go")
	cmd.Env = append(cmd.Env, "GOCACHE=C:\\Users\\1\\AppData\\Local\\go-build")
	cmd.Env = append(cmd.Env, "GOMODCACHE=C:\\Users\\1\\go\\pkg\\mod")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	cmd.Start()
	output, err := ioutil.ReadAll(stdout)
	errorOutput, err := ioutil.ReadAll(stderr)
	fmt.Println("Stdout:", string(output))
	fmt.Println("Stderr:", string(errorOutput))
	cmd.Wait()
}
func EncodeBinaryFile(pathCompile string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	defer encoder.Close()
	b, err := ioutil.ReadFile(pathCompile)
	if err != nil {
		return nil, err
	}
	encoder.Write(b)
	return &buf, nil
}
func DecodeBinaryFile(buffer *bytes.Buffer) ([]byte, error) {
	var dst []byte = make([]byte, base64.StdEncoding.DecodedLen(len(buffer.Bytes())))
	n, err := base64.StdEncoding.Decode(dst, buffer.Bytes())
	if err != nil {
		return []byte{}, err
	}
	return dst[:n], err
}
func CreateExecutableFile(src []byte) {
	file, err := os.Create("main.exe")
	if err != nil {
		fmt.Println(err)
	}
	file.Write(src)

}
func main() {
	const data = `package main
	func main(){
		fmt.Println(10)
	}
	`
	f, err := FormatFmt(data)
	if err != nil {
		fmt.Println(err)
	}
	ok, err := CheckCodePackageIsMain(f)
	if ok {
		fmt.Println("package is main")
	} else {
		fmt.Println("package not main")
	}
	p, o, err := WriteCodeFile(f)
	if err != nil {
		fmt.Println(err)
	}
	CompileCode(p, o)
	buf, err := EncodeBinaryFile(p)
	dst, err := DecodeBinaryFile(buf)
	CreateExecutableFile(dst)
	// CompileCode()

}
