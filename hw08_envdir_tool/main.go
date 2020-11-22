package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprint(os.Stdout, "you should provide dir and command [with args]\n")
		return
	}
	dirPath := os.Args[1]
	command := os.Args[2:]

	environments, err := ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	code := RunCmd(command, environments)

	os.Exit(code)
}
