package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	asciiart string = `
   ___                      __           
  / _ \_ __ ___  ___ _ __  / _|_ __ __ _ 
 / /_\/ '__/ _ \/ _ \ '_ \| |_| '__/ _` + "`" + ` | 
/ /_\\| | |  __/  __/ | | |  _| | | (_| |
\____/|_|  \___|\___|_| |_|_| |_|  \__,_|
                                         `
	command string
)

func init() {
	flag.Parse()
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
}

func main() {
	fmt.Println(asciiart)

	switch command {
	case "help":
		fmt.Println("No command for now.")
	default:
		if command != "" {
			fmt.Println("Unknown command:", command)
		}
		fmt.Println("Use 'go run main.go help' for usage information.")
	}
}
