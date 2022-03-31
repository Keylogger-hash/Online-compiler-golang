package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
)

type Manager struct {
}

func WriteCodeFile(text []byte) {
	currDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile(currDir+"/tmp/main.go", text, fs.FileMode(os.O_WRONLY))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	currDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < 1000; i++ {
		go func() {
			text := "package main\nimport \"fmt\"\nfunc main(){\n	fmt.Println(\"I am was writing\")\n}"
			WriteCodeFile([]byte(text))
			pathCompile := currDir + "/tmp/main.go"
			cmd := exec.Command("/opt/homebrew/bin/go", "run", pathCompile)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(output))
		}()

	}
	fmt.Scanln()
	fmt.Println("Process finished!")
}
