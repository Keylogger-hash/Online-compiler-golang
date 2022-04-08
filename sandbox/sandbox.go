package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
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
	Res   string
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

func (c *Container) StartContainer(decodeBytes []byte) ([]byte, error) {
	cmd := exec.Command("docker", "run", "-i", c.Name)
	var stdin,stdout, stderr *bytes.Buffer
	stdin.Write(decodeBytes)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		errorRunning := fmt.Errorf("%s\n%s",stderr.String(),err)
		return nil,errorRunning
	}
	output := stdout.Bytes()
	return output,nil
}
func DecodeBase64String(body string) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(body)))
	src := []byte(body)
	n, err := base64.StdEncoding.Decode(dst, src)
	if err != nil {
		return nil, err
	}
	return dst[:n], nil
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
		if err != nil {
			http.Error(wr, err.Error(), http.StatusBadRequest)
		}
		req := &Request{}
		resp := &Response{}

		err = json.Unmarshal(body, req)
		if err != nil {
			http.Error(wr, err.Error(), http.StatusBadRequest)
		}
		decodeBytes, err := DecodeBase64String(req.Body)
		if err != nil {
			resp.Error = err.Error()
			resp.Res = ""
			outputErr, _ := json.Marshal(resp)
			wr.Header().Add("Content-type", "application/json")
			wr.Write(outputErr)
		}
		c := &Container{Name: "sandbox-golang"}
		//_, err = c.buildContainer([]byte(req.Body))

		outputContainer, err := c.StartContainer(decodeBytes)
		if err != nil {
			resp.Error = err.Error()
			resp.Res = ""
			outputErr, _ := json.Marshal(resp)
			wr.Header().Add("Content-type", "application/json")
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleMain)
	mux.HandleFunc("/run", handleRun)
	http.ListenAndServe("localhost:8081", mux)
}
