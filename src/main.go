package main

import (
	"flag"
	"fmt"
)

var asciiart string = `
___                      __           
/ _ \_ __ ___  ___ _ __  / _|_ __ __ _ 
/ /_\/ '__/ _ \/ _ \ '_ \| |_| '__/ _` + "`" + ` | 
/ /_\\| | |  __/  __/ | | |  _| | | (_| |
\____/|_|  \___|\___|_| |_|_| |_|  \__,_|`

func init() {
	flag.Parse()
}

func main() {
	fmt.Printf(asciiart)
}
