package cmd

import (
	"fmt"
	"greenfra/src/services"
	"greenfra/src/utils"
	"log"

	"github.com/fatih/color"
)

func ListInstanceTypes(executePlan bool, planPath string) {
	if executePlan {
		err := utils.ExecuteTerraformPlan(planPath)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	plan, err := utils.ExecuteTerraformShow(planPath)
	if err != nil {
		log.Fatalf("%v", err)
	}

	region, err := utils.GetAWSRegion(plan)
	if err != nil {
		log.Fatalf("%v", err)
	}

	color.New(color.FgHiGreen).Printf("\nAWS Region: ")
	fmt.Println(region)

	// Extract resource changes
	changes, err := utils.ExtractResourceChanges(plan)
	if err != nil {
		log.Fatalf("Error extracting resource changes: %v", err)
	}

	// Analyze EC2 instances
	cfg := utils.LoadAWSConfig()
	ec2Service := services.NewEC2Service(cfg)
	err = ec2Service.Analyze(changes, region)
	if err != nil {
		log.Fatalf("Failed to analyze EC2: %v", err)
	}
	lambdaService := services.NewLambdaService(cfg)
	err = lambdaService.Analyze(changes, region)
	if err != nil {
		log.Fatalf("Failed to analyze Lambda: %v", err)
	}

}
