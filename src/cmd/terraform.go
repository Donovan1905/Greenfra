package cmd

import (
	"greenfra/src/services"
	"greenfra/src/utils"
	"log"
)

func HandleTerraform() (map[string]interface{}, error) {
	err := utils.ExecuteTerraformPlan()
	if err != nil {
		return nil, err
	}

	result, err := utils.ExecuteTerraformShow()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ListInstanceTypes() {
	// Get the Terraform plan JSON output
	plan, err := HandleTerraform()
	if err != nil {
		log.Fatalf("Error handling Terraform: %v", err)
	}

	// Extract resource changes
	changes, err := utils.ExtractResourceChanges(plan)
	if err != nil {
		log.Fatalf("Error extracting resource changes: %v", err)
	}

	// Analyze EC2 instances
	cfg := utils.LoadAWSConfig()
	ec2Service := services.NewEC2Service(cfg)
	err = ec2Service.Analyze(changes)
	if err != nil {
		log.Fatalf("Failed to analyze EC2: %v", err)
	}
}
