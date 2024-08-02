package cmd

import (
	"greenfra/src/services"
	"greenfra/src/utils"
	"log"
)

func HandleEC2() {
	cfg := utils.LoadAWSConfig()
	ec2Service := services.NewEC2Service(cfg)

	// Execute Terraform plan and show commands
	plan, err := HandleTerraform()
	if err != nil {
		log.Fatalf("Error handling Terraform: %v", err)
	}

	// Extract resource changes
	changes, err := utils.ExtractResourceChanges(plan)
	if err != nil {
		log.Fatalf("Error extracting resource changes: %v", err)
	}

	err = ec2Service.Analyze(changes)
	if err != nil {
		log.Fatalf("Failed to analyze EC2: %v", err)
	}
}
