package services

import (
	"fmt"
	"greenfra/src/types"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/olekukonko/tablewriter"
)

const defaultMeanLambdaExecutionDurationMilliseconds = 200
const defaultMeanNumberOfExecution = 100000

type LambdaService struct {
	client *lambda.Client
}

func NewLambdaService(cfg aws.Config) *LambdaService {
	return &LambdaService{
		client: lambda.NewFromConfig(cfg),
	}
}

func (s *LambdaService) Analyze(changes []ResourceChange, region string, comments map[string]types.ResourceMetadata) error {
	lambdasSpecs := make([]struct {
		name                                    string
		memorySize                              int
		vcpus                                   float64
		meanLambdaExecutionDurationMilliseconds int
		meanNumberOfExecution                   int
	}, 0)

	for _, change := range changes {
		if change.Type == "aws_lambda_function" {

			if after, ok := change.Change["after"].(map[string]interface{}); ok {
				name := change.Address

				memorySizeFloat, ok := after["memory_size"].(float64)
				if !ok {
					log.Printf("Warning: memory_size is not a float64 or is missing in change: %v", change.Change)
					continue
				}
				memorySize := int(memorySizeFloat)
				vcpus := float64(memorySize) / 1769

				var meanLambdaExecutionDurationMilliseconds int
				var meanNumberOfExecution int

				if resource, exists := comments[name]; exists {
					var err error
					meanLambdaExecutionDurationMilliseconds, err = strconv.Atoi(resource.Metadata["mean_execution_time"])
					if err != nil {
						log.Fatalf("mean_execution_time is not int : %v", err)
					}
					meanNumberOfExecution, err = strconv.Atoi(resource.Metadata["monthly_invocation"])
					if err != nil {
						log.Fatalf("monthly_invocation is not int : %v", err)
					}
				} else {
					meanLambdaExecutionDurationMilliseconds = defaultMeanLambdaExecutionDurationMilliseconds
					meanNumberOfExecution = defaultMeanNumberOfExecution
				}

				lambdasSpecs = append(lambdasSpecs, struct {
					name                                    string
					memorySize                              int
					vcpus                                   float64
					meanLambdaExecutionDurationMilliseconds int
					meanNumberOfExecution                   int
				}{name: name, memorySize: memorySize, vcpus: vcpus, meanLambdaExecutionDurationMilliseconds: meanLambdaExecutionDurationMilliseconds, meanNumberOfExecution: meanNumberOfExecution})
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
	name                                    string
	memorySize                              int
	vcpus                                   float64
	meanLambdaExecutionDurationMilliseconds int
	meanNumberOfExecution                   int
}, region string) {
	fmt.Print("\n")
	table := tablewriter.NewWriter(os.Stdout)
	headers := []string{"Lambda resource name", "vCPUs", "Memory (MiB)", "Estimated Monthly Power Consumption (kWh)", "Carbon impact (gCO2eq)"}

	for i, header := range headers {
		headers[i] = fmt.Sprintf("\x1b[32m%s\x1b[0m", header)
	}
	table.SetHeader(headers)

	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, spec := range lambdaSpecs {

		powerConsumption := calculateMonthlyPowerConsumption(
			spec.vcpus,
			spec.memorySize,
			(time.Duration(float64(spec.meanLambdaExecutionDurationMilliseconds))*time.Millisecond).Hours()*float64(spec.meanNumberOfExecution)) / 1000 // Convert Wh to kWh

		carbonImpact := calculateCarbonFootprint(powerConsumption, region)

		table.Append([]string{
			spec.name,
			fmt.Sprintf("%.1f", spec.vcpus),
			fmt.Sprintf("%d", spec.memorySize),
			fmt.Sprintf("%.10f", powerConsumption),
			fmt.Sprintf("%d", int(carbonImpact)),
		})
	}

	table.Render()
}
