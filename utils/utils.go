package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/tools/imports"
)

const maxBuildTime = time.Second * 5
const maxBinarySize = 10 << 20

type BuildResults struct {
	PathCompile string
	PathOutput  string
	Data        []byte
	Stdout      io.ReadCloser
	Stderr      io.ReadCloser
	Errors      error
	Events      []string
}

func NewBuildResult() *BuildResults {
	return &BuildResults{}
}
func CheckCodePackageIsMain(code []byte) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", code, parser.AllErrors)
	if err != nil {
		return err
	}
	if file.Name.Name != "main" {
		return fmt.Errorf("package is not main")
	}
	return nil
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
func WriteCodeFile(data []byte) (string, error) {
	uid := uuid.New().String()
	nameContainer := "main_" + uid
	PathOutput := fmt.Sprintf("/tmp/build/%s/", uid)
	PathCompile := PathOutput + nameContainer
	err := os.MkdirAll(PathOutput, 0750)
	if err != nil {
		return PathCompile, err
	}
	file, err := os.Create(PathCompile + ".go")
	if err != nil {
		return PathCompile, err
	}
	_, err = file.Write(data)
	if err != nil {
		return PathCompile, err
	}
	return PathCompile, nil
}

// func (b *BuildResults) cleanup(){
// 	os.RemoveAll(b.PathCompile)
// }
func EncodeToBase64(src []byte) []byte {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst
}
func CompileCode(tmpDir string) ([]byte, error) {
	defer os.RemoveAll("/tmp/build")

	var stderr strings.Builder
	ctx, finish := context.WithTimeout(context.Background(), maxBuildTime)
	defer finish()

	cmd := exec.CommandContext(ctx, "go", "build", "-o", tmpDir, tmpDir+".go")
	cmd.Env = append(cmd.Env, "GOOS=linux")
	cmd.Env = append(cmd.Env, "GOARCH=arm64")
	cmd.Env = append(cmd.Env, "GOPATH=/Users/pavelmorozov/go")
	cmd.Env = append(cmd.Env, "GOCACHE=/Users/pavelmorozov/Library/Caches/go-build")
	cmd.Stderr = &stderr

	err := cmd.Run()
	if ctx.Err() != nil {
		return nil, errors.New("Timeout is happened")
	}
	if err != nil {
		buildError := fmt.Errorf("%s", err.Error()+"\n"+stderr.String())
		return nil, buildError
	}
	file, _ := os.Stat(tmpDir)
	if buildSize := file.Size(); buildSize > maxBinarySize {
		maxBuildError := fmt.Errorf("Maximum binary file size exceeded")
		return nil, maxBuildError
	}
	data, err := ioutil.ReadFile(tmpDir)
	if err != nil {
		return nil,err
	}
	return data, nil

}

//
func EncodeBinaryFile(build *BuildResults) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	defer encoder.Close()
	b, err := ioutil.ReadFile(build.PathCompile)
	if err != nil {
		return nil, err
	}
	encoder.Write(b)
	return &buf, nil
}
