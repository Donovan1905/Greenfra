package main

import (
	"context"
	"flag"
	"fmt"
	"log"

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

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	switch command {
	case "ec2":
		fmt.Println(instanceType)
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
	case "help":
		fmt.Println("No command for now.")
	default:
		if command != "" {
			fmt.Println("Unknown command:", command)
		}
		fmt.Println("Use 'go run main.go help' for usage information.")
	}
}
