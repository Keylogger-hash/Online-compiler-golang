package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	pb "compiler.com/sandboxproto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	startTimeout  = time.Second * 30
	runTimeout    = time.Second * 5
	maxBinarySize = 100 << 20
	maxOutput     = 10 << 20
	memLimitBytes = 100 << 20
)

var httpServer *http.Server
var readyContainer chan *Container
var in chan []byte = make(chan []byte, 100)
var out chan []byte = make(chan []byte, 100)
var errC chan error = make(chan error, 100)

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
	cmd := exec.Command("sudo", "docker", "build", "-t", "sandbox-play", "-f", "play.Dockerfile", ".")
	err := cmd.Run()
	return err
}
func (c *Container) startContainer(decodeBytes []byte) ([]byte, error) {
	ctx, finish := context.WithTimeout(context.Background(), time.Second*30)
	defer finish()
	cmd := exec.CommandContext(ctx, "sudo", "docker", "run", "-i", "--rm", c.Name)
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
	if ctx.Err() != nil {
		errorTimeout := fmt.Errorf("Timelimit exceed")
		fmt.Println("Timeout exceed")
		cmd.Process.Kill()
		return nil, errorTimeout
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

type Server struct {
	pb.UnimplementedRunSandboxCompileCodeServer
}

func (s *Server) RunSandboxCompileCode(context context.Context, message *pb.RequestMessage) (*pb.ResponseMessage, error) {
	log.Printf("Run Sandbox Compile code")
	decodeBytes, err := DecodeBase64String(message.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(decodeBytes)
	in <- decodeBytes

	output, err := GetContainer(out, errC)
	fmt.Println(output)
	fmt.Println(err)
	if err != nil {
		return &pb.ResponseMessage{Res: string(output), Error: err.Error()}, nil
	}
	return &pb.ResponseMessage{Res: string(output), Error: ""}, nil
}

func GetContainer(out chan []byte, errC chan error) ([]byte, error) {
	select {
	case val := <-out:
		return val, nil
	case err := <-errC:
		return nil, err
	}
}
func workerLoop(in, out chan []byte, errC chan error) {
	select {
	case decodeBytes := <-in:
		c := Container{Name: "sandbox-play"}

		output, err := c.startContainer(decodeBytes)

		if err != nil {
			errC <- err
		}
		out <- output
	}

}
func makeWorkers() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go workerLoop(in, out, errC)
	}
}
func main() {
	ListContainers()
	err := buildContainer()
	if err != nil {
		fmt.Println(err)
		log.Fatal("Not build container sandbox-play.Stopped!")
	}
	fmt.Println("Build container")
	makeWorkers()
	l, err := net.Listen("tcp", "0.0.0.0:8082")
	if err != nil {
		log.Fatalf("Failed error to listen server port: %v", err)
	}
	fmt.Println("Listen  port :8082")
	server := grpc.NewServer()
	pb.RegisterRunSandboxCompileCodeServer(server, &Server{})
	err = server.Serve(l)
	fmt.Println("Serve port 8082")
	fmt.Println("Server starting...")
	if err != nil {
		log.Fatalf("Failed error to listen grpc server port: %v", err)
	}
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
