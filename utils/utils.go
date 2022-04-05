package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

func CompileCode() {
	
	cmd := exec.Command("go", "build", "tmp/main.go")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	cmd.Start()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if err := cmd.Start(); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := cmd.Wait(); err != nil {
	// 	log.Fatal(err)
	
	output, err := ioutil.ReadAll(stdout)
	errorOutput, err := ioutil.ReadAll(stderr)
	fmt.Println("Stdout:", string(output))
	fmt.Println("Stderr:", string(errorOutput))
	cmd.Wait()


}

func main() {
	CompileCode()

}
