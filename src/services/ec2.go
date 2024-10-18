package services

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/olekukonko/tablewriter"
)

type EC2Service struct {
	client *ec2.Client
}

func NewEC2Service(cfg aws.Config) *EC2Service {
	return &EC2Service{
		client: ec2.NewFromConfig(cfg),
	}
}

func (s *EC2Service) Analyze(changes []ResourceChange, region string) error {
	instanceTypes := make([]string, 0)

	for _, change := range changes {
		if change.Type == "aws_instance" {

			// Navigate into the 'after' map to get the 'instance_type'
			if after, ok := change.Change["after"].(map[string]interface{}); ok {
				instanceType, ok := after["instance_type"].(string)
				if !ok {
					log.Printf("Warning: instance_type is not a string or is missing in change: %v", change.Change)
					continue
				}
				instanceTypes = append(instanceTypes, instanceType)
			} else {
				log.Printf("Warning: 'after' map is missing in change: %v", change.Change)
			}
		}
	}

	if len(instanceTypes) > 0 {
		s.printInstanceSpecs(instanceTypes, region)
	}
	return nil
}

func (s *EC2Service) printInstanceSpecs(instanceTypes []string, region string) {
	awsInstanceTypes := make([]types.InstanceType, len(instanceTypes))
	for i, instanceType := range instanceTypes {
		awsInstanceTypes[i] = types.InstanceType(instanceType)
	}

	input := &ec2.DescribeInstanceTypesInput{
		InstanceTypes: awsInstanceTypes,
	}

	result, err := s.client.DescribeInstanceTypes(context.TODO(), input)
	if err != nil {
		log.Fatalf("Failed to describe instance types: %v", err)
	}
	fmt.Print("\n")
	// Create a table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Instance Type", "vCPUs", "Memory (MiB)", "Estimated Monthly Power Consumption (kWh)", "Carbon impact (gCO2eq)"})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.FgHiGreenColor}, tablewriter.Colors{tablewriter.FgHiGreenColor}, tablewriter.Colors{tablewriter.FgHiGreenColor}, tablewriter.Colors{tablewriter.FgHiGreenColor}, tablewriter.Colors{tablewriter.FgHiGreenColor})
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	// Add instance type details to the table
	for _, instanceType := range result.InstanceTypes {
		vcpus := *instanceType.VCpuInfo.DefaultVCpus
		memory := *instanceType.MemoryInfo.SizeInMiB

		// Calculate monthly power consumption
		powerConsumption := calculateMonthlyPowerConsumption(float64(vcpus), int(memory), hoursInMonth) / 1000

		// Calculate monthly carbon impact
		carbonImpact := calculateCarbonFootprint(powerConsumption, region)

		// Append the details including power consumption to the table
		table.Append([]string{
			string(instanceType.InstanceType),
			fmt.Sprintf("%d", vcpus),
			fmt.Sprintf("%d", memory),
			fmt.Sprintf("%.2f", powerConsumption),
			fmt.Sprintf("%d", int(carbonImpact)),
		})
	}

	// Render the table
	table.Render()
}
