package main

import (
	"fmt"
	"go/format"
	"golang.org/x/tools/imports"
)

func main() {
	var data string = `
	package main
	func main(){
		fmt.Println()
}
	`
	dest, err := format.Source([]byte(data))
	finish, err := imports.Process("",dest,nil)

	fmt.Println(string(finish), err)
}
