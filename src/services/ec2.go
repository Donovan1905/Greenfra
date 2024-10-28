package services

import (
	"context"
	"fmt"
	greenfraTypes "greenfra/src/types"
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

func (s *EC2Service) Analyze(changes []ResourceChange, region string, comments map[string]greenfraTypes.ResourceMetadata) error {
	instanceSpecs := make(map[string][]greenfraTypes.ResourceMetadata) // Map to hold instance types and their references

	for _, change := range changes {
		if change.Type == "aws_instance" {
			if after, ok := change.Change["after"].(map[string]interface{}); ok {
				instanceType, ok := after["instance_type"].(string)
				if !ok {
					log.Printf("Warning: instance_type is not a string or is missing in change: %v", change.Change)
					continue
				}
				// Capture the resource reference and add it to the map
				resourceMetadata := greenfraTypes.ResourceMetadata{
					ResourceReference: change.Address,
					// You can initialize Metadata as needed
				}
				instanceSpecs[instanceType] = append(instanceSpecs[instanceType], resourceMetadata)
			} else {
				log.Printf("Warning: 'after' map is missing in change: %v", change.Change)
			}
		}
	}

	if len(instanceSpecs) > 0 {
		s.printInstanceSpecs(instanceSpecs, region)
	}
	return nil
}

func (s *EC2Service) printInstanceSpecs(instanceSpecs map[string][]greenfraTypes.ResourceMetadata, region string) {
	awsInstanceTypes := make([]types.InstanceType, 0)

	for instanceType := range instanceSpecs {
		awsInstanceTypes = append(awsInstanceTypes, types.InstanceType(instanceType))
	}

	input := &ec2.DescribeInstanceTypesInput{
		InstanceTypes: awsInstanceTypes,
	}

	result, err := s.client.DescribeInstanceTypes(context.TODO(), input)
	if err != nil {
		log.Fatalf("Failed to describe instance types: %v", err)
	}
	fmt.Print("\n")
	table := tablewriter.NewWriter(os.Stdout)
	headers := []string{"Instance Type", "Resource References", "vCPUs", "Memory (MiB)", "Estimated Monthly Power Consumption (kWh)", "Carbon impact (gCO2eq)"}
	table.SetHeader(headers)

	colors := make([]tablewriter.Colors, len(headers))
	for i := range colors {
		colors[i] = tablewriter.Colors{tablewriter.FgHiGreenColor}
	}
	table.SetHeaderColor(colors...)

	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, instanceType := range result.InstanceTypes {
		vcpus := *instanceType.VCpuInfo.DefaultVCpus
		memory := *instanceType.MemoryInfo.SizeInMiB

		powerConsumption := calculateMonthlyPowerConsumption(float64(vcpus), int(memory), hoursInMonth) / 1000
		carbonImpact := calculateCarbonFootprint(powerConsumption, region)

		references := instanceSpecs[string(instanceType.InstanceType)]
		referenceList := ""
		for _, ref := range references {
			referenceList += ref.ResourceReference + ", "
		}
		if len(referenceList) > 0 {
			referenceList = referenceList[:len(referenceList)-2]
		}

		table.Append([]string{
			string(instanceType.InstanceType),
			referenceList,
			fmt.Sprintf("%d", vcpus),
			fmt.Sprintf("%d", memory),
			fmt.Sprintf("%.2f", powerConsumption),
			fmt.Sprintf("%d", int(carbonImpact)),
		})
	}

	table.Render()
}
