package main

import (
	"encoding/json"
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

type Request struct {
	Body string
}
type Response struct {
	Res string
	Error string
}
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
	defer file.Close()
	cmd := exec.Command("sudo","docker", "build", ".", "-t", nameContainer, "--build-arg", argName)
	fmt.Println(cmd)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(string(output))
	}
	fmt.Println(string(output))
	return output, nil
}
func (c *Container) StartContainer() ([]byte, error) {
	cmd := exec.Command("sudo","docker","run","-i",c.Name)
	fmt.Println(cmd)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
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
		body, err := ioutil.ReadAll(r.Body)
		fmt.Println(string(body))
		if err != nil {
			http.Error(wr, err.Error(), http.StatusBadRequest)
		}
		req := &Request{}
		err = json.Unmarshal(body, req)
		fmt.Println(err)
		if err != nil {
			http.Error(wr, err.Error(),http.StatusBadRequest)
		}
		fmt.Println(string(req.Body))
		c :=&Container{}
		resp := &Response{}
		_,err =  c.buildContainer([]byte(req.Body))
		if err != nil {
			resp.Error = err.Error()
			resp.Res = ""
			output, _ := json.Marshal(resp)
			wr.Header().Add("Content-type","application/json")
			wr.Write(output)
		}
		outputContainer, err := c.StartContainer()
		if err != nil {
			resp.Error = err.Error()
			resp.Res = ""
			outputErr, _ := json.Marshal(resp)
			wr.Header().Add("Content-type","application/json")
			wr.Write(outputErr)
		}
		resp.Res = string(outputContainer)
		resp.Error = ""
		output, _ := json.Marshal(resp)
		fmt.Println(string(outputContainer))
		wr.Write(output)
		
	default:
		wr.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func main() {
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", handleMain)
	// mux.HandleFunc("/run", handleRun)
	// http.ListenAndServe("localhost:8081", mux)
 }
