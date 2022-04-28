package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"
	"time"
)

type limitedWriterByte struct {
	limit int
	out   *bytes.Buffer
}

func startContainer(ctx context.Context, decodeBytes []byte, cancel context.CancelFunc) ([]byte, error) {
	// ctx, finish := context.WithTimeout(context.Background(), time.Second*5)
	// defer finish()
	cmd := exec.CommandContext(ctx, "sudo", "docker", "run", "-i", "--rm", "sandbox-play")
	var stdin, stdout, stderr bytes.Buffer
	defer cancel()

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

var readyContainer chan []byte
var errC chan error

func GetContainer(ctx context.Context) ([]byte, error) {
	select {
	case val := <-readyContainer:
		return val, nil
	case err := <-errC:
		return nil, err
	case <-ctx.Done():
		err := fmt.Errorf("Timeout is reached")
		return nil, err
	}
}
func worker(ctx context.Context, finish context.CancelFunc) {

	data, _ := ioutil.ReadFile("main1")
	c, err := startContainer(ctx, data, finish)
	if err != nil {
		errC <- err
	}
	readyContainer <- c

}
func main() {
	for i := 0; i < runtime.NumCPU(); i++ {
		data, _ := ioutil.ReadFile("main1")
		ctx, finish := context.WithTimeout(context.Background(), time.Microsecond)
		c, err := startContainer(ctx, data, finish)
		fmt.Println(err)
		if c == nil {
			fmt.Println()
		} else {
			fmt.Println(string(c))

		}
		// ctx, finish := context.WithTimeout(context.Background(), time.Second*5)
		// go worker(ctx, finish)
		// d, err := GetContainer(ctx)
		// fmt.Println(string(d))
		// fmt.Println(err)
	}

	fmt.Scanln()
}
