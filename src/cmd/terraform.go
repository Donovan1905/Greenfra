package cmd

import (
	"greenfra/src/services"
	"greenfra/src/utils"
	"log"
)

// func HandleTerraform() (map[string]interface{}, error) {

// 	result, err := utils.ExecuteTerraformShow()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

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
