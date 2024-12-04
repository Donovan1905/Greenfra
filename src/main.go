package main

import (
	"flag"
	"fmt"
	"os"

	"greenfra/src/cmd"
)

var (
	command     string
	planPath    string = "greenfra.tfplan"
	executePlan bool
)

func init() {
	flag.BoolVar(&executePlan, "exec-plan", false, "Specify whether or not greenfra should execute terraform plan or you provide the tfplan file")

	if envValue, exists := os.LookupEnv("GREENFRA_EXEC_PLAN"); exists && envValue == "true" {
		executePlan = true
	}

	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		command = args[0]
	}
	if len(args) >= 2 {
		planPath = args[1]
	}
}

func main() {
	fmt.Printf("\x1b[32m%s\x1b[0m\n", cmd.AsciiArt)

	switch command {
	case "help":
		fmt.Println("Usage: go run main.go [command] [flags]")
		fmt.Println("Commands:")
		fmt.Println("  analyze <tfplan file path>      - List all the instances in your terraform plan file")
		fmt.Println("Flags:")
		flag.PrintDefaults()
	case "analyze":
		cmd.ListResources(executePlan, planPath)
	default:
		fmt.Println("Unknown command")
	}
}
