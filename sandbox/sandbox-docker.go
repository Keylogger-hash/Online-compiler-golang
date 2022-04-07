package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func RunCode(stdin []byte) {
	file, err := os.Create("binary")
	if err != nil {
		fmt.Println(err)
	}
	file.Write(stdin)
	file.Close()
	cmd := exec.Command("./binary")
	stdout, err := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		fmt.Println(err)
	}
	b, err := ioutil.ReadAll(stdout)
	fmt.Println(string(b))
	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
	}
}
func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err)
	}
	RunCode(b)

}
