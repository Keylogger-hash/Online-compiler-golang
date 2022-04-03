package main

import (
	"fmt"
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

type Sandbox struct {
}

type Container struct {
}

func buildCode() {

}
func ListContainers() {
	cmd := exec.Command("sudo", "docker", "ps", "--format", "'{{json .}}'", "-a")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(output))
}
func buildContainer(){
	nameContainer  := "main_"+time.Now().String()
	cmd := exec.Command("sudo","docker","build","-it",nameContainer,".")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(output)
	}
	fmt.Println(string(output))
}
func StartContainer(){
	
}
func RunCode() {
	cmd := exec.Command("./main")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}
func main() {
	buildContainer()
}
