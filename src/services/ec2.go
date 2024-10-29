package services

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/olekukonko/tablewriter"
	greenfraTypes "greenfra/src/types"
	"log"
	"os"
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
	ec2Specs := make([]struct {
		name         string
		instanceType string
	}, 0)

	for _, change := range changes {
		if change.Type == "aws_instance" {

			if after, ok := change.Change["after"].(map[string]interface{}); ok {
				name := change.Address

				instanceType, ok := after["instance_type"].(string)
				if !ok {
					log.Printf("Warning: instance type is not a string or is missing in change: %v", change.Change)
					continue
				}

				ec2Specs = append(ec2Specs, struct {
					name         string
					instanceType string
				}{name: name, instanceType: instanceType})
			} else {
				log.Printf("Warning: 'after' map is missing in change: %v", change.Change)
			}
		}
	}

	if len(ec2Specs) > 0 {
		s.printInstanceSpecs(ec2Specs, region)
	}
	return nil
}

func (s *EC2Service) printInstanceSpecs(ec2Specs []struct {
	name         string
	instanceType string
}, region string) {

	instanceTypeMap := make(map[string]struct{})
	for _, spec := range ec2Specs {
		instanceTypeMap[spec.instanceType] = struct{}{}
	}

	awsInstanceTypes := make([]types.InstanceType, 0)
	for instanceType := range instanceTypeMap {
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
	headers := []string{"Resource Reference", "Instance Type", "vCPUs", "Memory (MiB)", "Estimated Monthly Power Consumption (kWh)", "Carbon impact (gCO2eq)"}
	table.SetHeader(headers)

	colors := make([]tablewriter.Colors, len(headers))
	for i := range colors {
		colors[i] = tablewriter.Colors{tablewriter.FgHiGreenColor}
	}
	table.SetHeaderColor(colors...)

	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	instanceTypeDetails := make(map[string]types.InstanceTypeInfo)

	for _, instanceType := range result.InstanceTypes {
		instanceTypeDetails[string(instanceType.InstanceType)] = instanceType
	}

	for _, spec := range ec2Specs {
		instanceTypeDetail, exists := instanceTypeDetails[spec.instanceType]
		if !exists {
			log.Printf("Warning: Instance type %s not found in DescribeInstanceTypes result.", spec.instanceType)
			continue
		}

		vcpus := *instanceTypeDetail.VCpuInfo.DefaultVCpus
		memory := *instanceTypeDetail.MemoryInfo.SizeInMiB

		powerConsumption := calculateMonthlyPowerConsumption(float64(vcpus), int(memory), hoursInMonth) / 1000
		carbonImpact := calculateCarbonFootprint(powerConsumption, region)

		table.Append([]string{
			spec.name,
			spec.instanceType,
			fmt.Sprintf("%d", vcpus),
			fmt.Sprintf("%d", memory),
			fmt.Sprintf("%.2f", powerConsumption),
			fmt.Sprintf("%d", int(carbonImpact)),
		})
	}

	table.Render()
}
