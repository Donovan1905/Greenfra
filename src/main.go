package main

import (
	"flag"
	"fmt"

	"greenfra/src/cmd"

	"github.com/fatih/color"
)

var (
	command     string
	planPath    string = "greenfra.tfplan"
	executePlan bool
)

func init() {
	flag.BoolVar(&executePlan, "exec-plan", false, "Specify wheter or not greenfra should execute terraform plan or you provide the tfplan file")
	flag.Parse()
	command = flag.Args()[0]
	if len(flag.Args()) >= 2 {
		planPath = flag.Args()[1]
	}
}

func main() {
	color.New(color.FgHiGreen).Printf(cmd.AsciiArt)

	switch command {
	case "help":
		fmt.Println("Usage: go run main.go [command] [flags]")
		fmt.Println("Commands:")
		fmt.Println("  analyze <tfplan file path>      - List all the instances in your terraform plan file")
		fmt.Println("Flags:")
		flag.PrintDefaults()
	case "analyze":
		cmd.ListInstanceTypes(executePlan, planPath)
	default:
		fmt.Println("Unkown command")
	}
}
