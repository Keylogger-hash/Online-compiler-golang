package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func RunCode(stdin []byte) {
	file, err := os.Create("binary")
	if err != nil {
		log.Fatal(err)
	}
	file.Write(stdin)
	err = file.Chmod(0777)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	cmd := exec.Command("./binary")
	stdout, err := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(stdout)
	fmt.Println(string(b))
	if err := cmd.Wait(); err != nil {
		os.Remove("binary")
		log.Fatal(err)
	}
}
func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	RunCode(b)
}
