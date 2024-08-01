package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

var (
	asciiart string = `
   ___                      __           
  / _ \_ __ ___  ___ _ __  / _|_ __ __ _ 
 / /_\/ '__/ _ \/ _ \ '_ \| |_| '__/ _` + "`" + ` | 
/ /_\\| | |  __/  __/ | | |  _| | | (_| |
\____/|_|  \___|\___|_| |_|_| |_|  \__,_|
                                         `
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
	fmt.Println(asciiart)

	switch command {
	case "ec2":
		handleEC2()
	case "terraform":
		handleTerraform()
	case "help":
		fmt.Println("No command for now.")
	default:
		if command != "" {
			fmt.Println("Unknown command:", command)
		}
		fmt.Println("Use 'go run main.go help' for usage information.")
	}
}

func handleEC2() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeInstanceTypesInput{
		InstanceTypes: []types.InstanceType{types.InstanceType(instanceType)},
	}

	result, err := svc.DescribeInstanceTypes(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to describe instance types, %v", err)
	}

	for _, instanceType := range result.InstanceTypes {
		vcpus := *instanceType.VCpuInfo.DefaultVCpus
		memory := *instanceType.MemoryInfo.SizeInMiB
		fmt.Printf("Instance Type: %s\n", instanceType.InstanceType)
		fmt.Printf("vCPUs: %d\n", vcpus)
		fmt.Printf("Memory (MiB): %d\n", memory)
	}
}

func handleTerraform() {
	// Execute "terraform plan -out tfplan"
	cmdPlan := exec.Command("terraform", "plan", "-out=tfplan")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmdPlan.Stdout = &out
	cmdPlan.Stderr = &stderr
	err := cmdPlan.Run()
	if err != nil {
		log.Fatalf("terraform plan failed: %v\n%s", err, stderr.String())
	}
	fmt.Println("Terraform plan executed successfully.")

	// Execute "terraform show -json tfplan"
	cmdShow := exec.Command("terraform", "show", "-json", "tfplan")
	var outShow bytes.Buffer
	var stderrShow bytes.Buffer
	cmdShow.Stdout = &outShow
	cmdShow.Stderr = &stderrShow
	err = cmdShow.Run()
	if err != nil {
		log.Fatalf("terraform show failed: %v\n%s", err, stderrShow.String())
	}

	// Parse the JSON output
	var result map[string]interface{}
	err = json.Unmarshal(outShow.Bytes(), &result)
	if err != nil {
		log.Fatalf("failed to parse JSON: %v", err)
	}

	// Filter to only keep the EC2 instances and retrieve their instance types
	if plannedValues, ok := result["planned_values"].(map[string]interface{}); ok {
		if rootModule, ok := plannedValues["root_module"].(map[string]interface{}); ok {
			if resources, ok := rootModule["resources"].([]interface{}); ok {
				for _, resource := range resources {
					if resMap, ok := resource.(map[string]interface{}); ok {
						if resType, ok := resMap["type"].(string); ok && resType == "aws_instance" {
							if values, ok := resMap["values"].(map[string]interface{}); ok {
								if instanceType, ok := values["instance_type"].(string); ok {
									fmt.Printf("Instance Type: %s\n", instanceType)
								}
							}
						}
					}
				}
			}
		}
	}
}
