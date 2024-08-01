package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/tabwriter"

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
		handleEC2(instanceType)
	case "terraform":
		_, err := handleTerraform()
		if err != nil {
			log.Fatalf("Error handling Terraform: %v", err)
		}
	case "help":
		fmt.Println("Usage: go run main.go [command] [flags]")
		fmt.Println("Commands:")
		fmt.Println("  ec2   - Describe EC2 instance types")
		fmt.Println("  terraform - Execute Terraform commands")
		fmt.Println("Flags:")
		flag.PrintDefaults()
	default:
		listInstanceTypes()
	}
}

func handleEC2(instanceType string) {
	// Load the shared AWS configuration (from ~/.aws/config or environment variables)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	if instanceType == "" {
		fmt.Println("Instance type must be specified with -instance-type flag.")
		return
	}

	// Create an EC2 service client
	svc := ec2.NewFromConfig(cfg)

	// Describe instance types
	input := &ec2.DescribeInstanceTypesInput{
		InstanceTypes: []types.InstanceType{types.InstanceType(instanceType)},
	}

	result, err := svc.DescribeInstanceTypes(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to describe instance types, %v", err)
	}

	// Create a tab writer
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

	// Print the headers
	fmt.Fprintln(w, "Instance Type\tvCPUs\tMemory (MiB)")

	// Print the details of the instance type with dereferenced values
	for _, instanceType := range result.InstanceTypes {
		vcpus := *instanceType.VCpuInfo.DefaultVCpus
		memory := *instanceType.MemoryInfo.SizeInMiB
		fmt.Fprintf(w, "%s\t%d\t%d\n", instanceType.InstanceType, vcpus, memory)
	}

	// Flush the writer
	w.Flush()
}

func handleTerraform() (map[string]interface{}, error) {
	// Execute "terraform plan -out tfplan"
	cmdPlan := exec.Command("terraform", "plan", "-out=tfplan")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmdPlan.Stdout = &out
	cmdPlan.Stderr = &stderr
	err := cmdPlan.Run()
	if err != nil {
		return nil, fmt.Errorf("terraform plan failed: %v\n%s", err, stderr.String())
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
		return nil, fmt.Errorf("terraform show failed: %v\n%s", err, stderrShow.String())
	}

	// Parse the JSON output
	var result map[string]interface{}
	err = json.Unmarshal(outShow.Bytes(), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return result, nil
}

func listInstanceTypes() {
	// Get the Terraform plan JSON output
	result, err := handleTerraform()
	if err != nil {
		log.Fatalf("Error handling Terraform: %v", err)
	}

	// Collect instance types
	instanceTypes := make(map[string]bool)

	if plannedValues, ok := result["planned_values"].(map[string]interface{}); ok {
		if rootModule, ok := plannedValues["root_module"].(map[string]interface{}); ok {
			if resources, ok := rootModule["resources"].([]interface{}); ok {
				for _, resource := range resources {
					if resMap, ok := resource.(map[string]interface{}); ok {
						if resType, ok := resMap["type"].(string); ok && resType == "aws_instance" {
							if values, ok := resMap["values"].(map[string]interface{}); ok {
								if instanceType, ok := values["instance_type"].(string); ok {
									instanceTypes[instanceType] = true
								}
							}
						}
					}
				}
			}
		}
	}

	// Create a tab writer
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

	// Print the headers
	fmt.Fprintln(w, "Instance Type\tvCPUs\tMemory (MiB)")

	// Retrieve and print specifications for each instance type from AWS
	if len(instanceTypes) > 0 {
		for instanceType := range instanceTypes {
			// Load the shared AWS configuration (from ~/.aws/config or environment variables)
			cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
			if err != nil {
				log.Fatalf("unable to load SDK config, %v", err)
			}

			// Create an EC2 service client
			svc := ec2.NewFromConfig(cfg)

			// Describe instance types
			input := &ec2.DescribeInstanceTypesInput{
				InstanceTypes: []types.InstanceType{types.InstanceType(instanceType)},
			}

			result, err := svc.DescribeInstanceTypes(context.TODO(), input)
			if err != nil {
				log.Fatalf("failed to describe instance types, %v", err)
			}

			// Print the details of the instance type with dereferenced values
			for _, instanceType := range result.InstanceTypes {
				vcpus := *instanceType.VCpuInfo.DefaultVCpus
				memory := *instanceType.MemoryInfo.SizeInMiB
				fmt.Fprintf(w, "%s\t%d\t%d\n", instanceType.InstanceType, vcpus, memory)
			}
		}
	} else {
		fmt.Println("No EC2 instances found in the Terraform plan.")
	}

	// Flush the writer
	w.Flush()
}
