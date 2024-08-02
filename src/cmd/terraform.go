package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	"greenfra/src/utils"
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
	result, err := HandleTerraform()
	if err != nil {
		log.Fatalf("Error handling Terraform: %v", err)
	}

	// Collect instance types
	instanceTypes := make(map[string]bool)

	if plannedValues, ok := result["planned_values"].(map[string]interface{}); ok {
		if rootModule, ok := plannedValues["root_module"].(map[string]interface{}); ok {
			if resources, ok := rootModule["resources"].([]interface{}); ok {
				for _, resource := range resources {
					if resMap, ok := resource.(map[string]interface{}); ok {
						if resType, ok := resMap["type"].(string); ok && resType == "aws_instance" {
							if values, ok := resMap["values"].(map[string]interface{}); ok {
								if instanceType, ok := values["instance_type"].(string); ok {
									instanceTypes[instanceType] = true
								}
							}
						}
					}
				}
			}
		}
	}

	// Create a table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Instance Type", "vCPUs", "Memory (MiB)"})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.FgHiGreenColor}, tablewriter.Colors{tablewriter.FgHiGreenColor}, tablewriter.Colors{tablewriter.FgHiGreenColor})
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	// Retrieve and add specifications for each instance type to the table
	if len(instanceTypes) > 0 {
		cfg := utils.LoadAWSConfig()
		svc := utils.CreateEC2Client(cfg)
		instanceTypesList := make([]string, 0, len(instanceTypes))
		for instanceType := range instanceTypes {
			instanceTypesList = append(instanceTypesList, instanceType)
		}

		instancesInfo := utils.DescribeInstanceTypes(svc, instanceTypesList)
		for _, info := range instancesInfo {
			table.Append([]string{info.InstanceType, fmt.Sprintf("%d", info.VCPUs), fmt.Sprintf("%d", info.Memory)})
		}
	} else {
		fmt.Println("No EC2 instances found in the Terraform plan.")
	}

	// Render the table
	table.Render()
}
