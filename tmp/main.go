package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
	"time"
)

type limitedWriterByte struct {
	limit int
	out   *bytes.Buffer
}

func startContainer(decodeBytes []byte) ([]byte, error) {
	ctx, finish := context.WithTimeout(context.Background(), time.Second*30)
	defer finish()
	cmd := exec.CommandContext(ctx, "sudo", "docker", "run", "-i", "--rm", "sandbox-play")
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
		cmd.Process.Kill()
		return nil, errorTimeout
	}

	output := stdout.Bytes()
	return output, nil
}

type Cmd struct {
	cmd *exec.Cmd
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
	case data := <-in:
		c, err := startContainer(data)
		if err != nil {
			errC <- err
		}
		out <- c
	}

}

var in chan []byte = make(chan []byte)
var out chan []byte = make(chan []byte)
var errC chan error = make(chan error)

func makeWorkers() {
	for i := 0; i < 10; i++ {
		go workerLoop(in, out, errC)
	}
}

// func makeWorkers() {
// 	for i := 0; i < runtime.NumCPU(); i++ {
// 		//ctx, finish := context.WithTimeout(context.Background(), time.Second*30)
// 		go workerLoop()

// 	}
// }
func main() {
	makeWorkers()

	data, _ := ioutil.ReadFile("main1")
	in <- data
	c, err := GetContainer(out, errC)
	fmt.Println(string(c))
	fmt.Println(err)
	// for i := 0; i < 10; i++ {
	// 	data, _ := ioutil.ReadFile("main1")
	// 	in <- data
	// }

	// for i:=0;i<10;i++{
	// 	c, err := GetContainer(out,errC)
	// 	fmt.Println(string(c))
	// 	fmt.Println(err)
	// }
	fmt.Scanln()
}
