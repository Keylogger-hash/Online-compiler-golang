package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"
)

const (
	startTimeout  = time.Second * 30
	runTimeout    = time.Second * 5
	maxBinarySize = 100 << 20
	maxOutput     = 100 << 20
	memLimitBytes = 100 << 20
)

var httpServer *http.Server

type Container struct {
	Name string

	Stdout []byte
	Stderr []byte
	Stdin  []byte

	Cmd *exec.Cmd
}

func ListContainers() {
	cmd := exec.Command("sudo", "docker", "ps", "--format", "'{{json .}}'", "-a")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(output))
}
func (c *Container) buildContainer(data []byte) ([]byte, error) {
	uid := uuid.New().String()
	nameContainer := "main_" + uid
	fileName := nameContainer + ".go"
	c.Name = nameContainer
	argName := fmt.Sprintf("NameFile=%s", nameContainer)
	file, err := os.Create("tmp/" + fileName)
	file.Write(data)
	cmd := exec.Command("sudo", "docker", "build", ".", "-t", nameContainer, "--build-arg", argName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
	}
	fmt.Println(string(output))
	return output, nil
}
func (c *Container) StartContainer() ([]byte, error) {
	cmdStr := fmt.Sprintf("docker run -i %s", c.Name)
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	fmt.Println(cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil,err
	}
	return output, nil
}

func handleMain(wr http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(wr, r)
		return
	}
	wr.Write([]byte("Everthing okay!"))
}
func handleRun(wr http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":

	default:
		wr.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func main() {
	c := &Container{}
	file, err := ioutil.ReadFile("tmp/main.go")
	if err != nil {
		fmt.Println(err)
	}
	c.buildContainer(file)
	fmt.Println("Container build")
	data, err := c.StartContainer()
	
	fmt.Println(string(data))
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", handleMain)
	// mux.HandleFunc("/run",handleRun)
	// http.ListenAndServe("localhost:8081", mux)
}
