package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/fatih/color"
	"greenfra/src/cmd"
)

var (
	command      string
	instanceType string
)

func init() {
	flag.StringVar(&instanceType, "instance-type", "", "Specify the EC2 instance type")
	flag.Parse()
	if len(flag.Args()) >= 1 {
		command = flag.Args()[0]
	}
}

func main() {
	color.New(color.FgHiGreen).Printf(cmd.AsciiArt)

	switch command {
	case "ec2":
		cmd.HandleEC2()
	case "terraform":
		_, err := cmd.HandleTerraform()
		if err != nil {
			log.Fatalf("Error handling Terraform: %v", err)
		}
	case "help":
		fmt.Println("Usage: go run main.go [command] [flags]")
		fmt.Println("Commands:")
		fmt.Println("  ec2        - Describe EC2 instance types")
		fmt.Println("  terraform  - Execute Terraform commands")
		fmt.Println("Flags:")
		flag.PrintDefaults()
	default:
		cmd.ListInstanceTypes()
	}
}
