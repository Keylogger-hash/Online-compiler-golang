package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"time"

	pb "compiler.com/sandboxproto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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
func buildContainer() error {
	cmd := exec.Command("sudo", "docker", "-f", "play.Dockerfile", "-t", "sandbox-play", ".")
	err := cmd.Run()
	return err
}
func (c *Container) startContainer(decodeBytes []byte) ([]byte, error) {
	cmd := exec.Command("sudo", "docker", "run", "-i", c.Name)
	var stdin, stdout, stderr bytes.Buffer
	stdin.Write(decodeBytes)
	cmd.Stdin = &stdin
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		errorRunning := fmt.Errorf("%s\n%s", stderr.String(), err)
		return nil, errorRunning
	}
	output := stdout.Bytes()
	return output, nil
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
func EncodeToBase64(src []byte) string {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return string(dst)
}

func handleMain(wr http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(wr, r)
		return
	}
	wr.Write([]byte("Everthing okay!"))
}

// func handleRun(wr http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case "POST":
// 		body, err := ioutil.ReadAll(r.Body)
// 		if err != nil {
// 			http.Error(wr, err.Error(), http.StatusBadRequest)
// 		}
// 		req := &Request{}
// 		resp := &Response{}

// 		err = json.Unmarshal(body, req)
// 		if err != nil {
// 			http.Error(wr, err.Error(), http.StatusBadRequest)
// 		}
// 		decodeBytes, err := DecodeBase64String(req.Body)
// 		if err != nil {
// 			resp.Error = err.Error()
// 			resp.Res = ""
// 			outputErr, _ := json.Marshal(resp)
// 			wr.Header().Add("Content-type", "application/json")
// 			wr.Write(outputErr)
// 		}
// 		c := &Container{Name: "sandbox-play"}
// 		//_, err = c.buildContainer([]byte(req.Body))

// 		outputContainer, err := c.startContainer(decodeBytes)
// 		if err != nil {
// 			resp.Error = err.Error()
// 			resp.Res = ""
// 			outputErr, _ := json.Marshal(resp)
// 			wr.Header().Add("Content-type", "application/json")
// 			wr.Write(outputErr)
// 		}
// 		resp.Res = string(outputContainer)
// 		resp.Error = ""
// 		output, _ := json.Marshal(resp)
// 		fmt.Println(string(outputContainer))
// 		wr.Write(output)

// 	default:
// 		wr.WriteHeader(http.StatusMethodNotAllowed)
// 	}
// }
type Server struct {
}

func (s *Server) RunSandboxCompileCode(context context.Context, message *pb.RequestMessage) (*pb.ResponseMessage, error) {
	log.Printf("Run Sandbox Compile code")
	return &pb.ResponseMessage{Res: "Hello world", Error: ""}, nil
}
func main() {
	l, err := net.Listen("tcp", "0.0.0.0:8081")
	if err != nil {
		log.Fatal("Failed error to listen server port: %v", err)
	}
	server := grpc.NewServer()
	rc := RunSandboxCompileCodeServer.Server{}
	pb.RegisterRunSandboxCompileCodeServer(server, &rc)
	err = server.Serve(l)
	if err != nil {
		log.Fatal("Failed error to listen grpc server port: %v", err)
	}
	fmt.Println("Server starting...")
	// ListContainers()
	// err := buildContainer()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// // c.Name = "8fbcde8f75bf"
	// // file, err := os.Open("main_ed723086-f775-42e5-8cc6-5a063a0bb330")
	// // defer file.Close()
	// // if err != nil {
	// // 	fmt.Println(err)
	// // }
	// // fileBytes, err := ioutil.ReadAll(file)
	// // if err != nil {
	// // 	fmt.Println(err)
	// // }
	// // dst := EncodeToBase64(fileBytes)
	// // fmt.Println(dst)
	// // b, err := DecodeBase64String(dst)
	// // if err != nil {
	// // 	fmt.Println(err)
	// // }
	// // out, err := c.StartContainer(b)
	// // if err != nil {
	// // 	fmt.Println(err)
	// // }
	// // fmt.Println(string(out))
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", handleMain)
	// mux.HandleFunc("/run", handleRun)
	// fmt.Println("Starting server...")
	// fmt.Println("Listen and serve on localhost 8081")
	// http.ListenAndServe("localhost:8081", mux)
}
