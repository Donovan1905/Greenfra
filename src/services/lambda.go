package services

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/olekukonko/tablewriter"
)

type LambdaService struct {
	client *lambda.Client
}

func NewLambdaService(cfg aws.Config) *LambdaService {
	return &LambdaService{
		client: lambda.NewFromConfig(cfg),
	}
}

func (s *LambdaService) Analyze(changes []ResourceChange, region string) error {
	lambdasSpecs := make([]struct {
		name       string
		memorySize int
	}, 0)

	for _, change := range changes {
		if change.Type == "aws_lambda_function" {

			if after, ok := change.Change["after"].(map[string]interface{}); ok {
				// Get the resource reference name
				name := change.Address

				// Change here: assert as float64 and then convert to int
				memorySizeFloat, ok := after["memory_size"].(float64)
				if !ok {
					log.Printf("Warning: memory_size is not a float64 or is missing in change: %v", change.Change)
					continue
				}
				memorySize := int(memorySizeFloat) // Convert to int

				// Append to lambdasSpecs with name and memory size
				lambdasSpecs = append(lambdasSpecs, struct {
					name       string
					memorySize int
				}{name: name, memorySize: memorySize})
			} else {
				log.Printf("Warning: 'after' map is missing in change: %v", change.Change)
			}
		}
	}

	if len(lambdasSpecs) > 0 {
		s.printLambdasSpecs(lambdasSpecs, region)
	}
	return nil
}

func (s *LambdaService) printLambdasSpecs(lambdaSpecs []struct {
	name       string
	memorySize int
}, region string) {
	fmt.Print("\n")
	// Create a table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Lambda resource name", "vCPUs", "Memory (MiB)", "Estimated Monthly Power Consumption (kWh)", "Carbon impact (gCO2eq)"})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.FgHiGreenColor}, tablewriter.Colors{tablewriter.FgHiGreenColor}, tablewriter.Colors{tablewriter.FgHiGreenColor}, tablewriter.Colors{tablewriter.FgHiGreenColor}, tablewriter.Colors{tablewriter.FgHiGreenColor})
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, spec := range lambdaSpecs {
		vcpus := float64(spec.memorySize) / 1769 // Assuming 1769 MB per vCPU

		// Calculate monthly power consumption
		powerConsumption := calculateMonthlyPowerConsumption(vcpus, spec.memorySize, hoursInMonth) / 1000 // Convert Wh to kWh

		// Calculate monthly carbon impact
		carbonImpact := calculateCarbonFootprint(powerConsumption, region)

		// Append the details including power consumption to the table
		table.Append([]string{
			spec.name,
			fmt.Sprintf("%.1f", vcpus),
			fmt.Sprintf("%d", spec.memorySize),
			fmt.Sprintf("%.2f", powerConsumption),
			fmt.Sprintf("%d", int(carbonImpact)),
		})
	}

	table.Render()
}
