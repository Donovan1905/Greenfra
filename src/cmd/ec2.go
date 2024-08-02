package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"greenfra/src/utils"
)

func HandleEC2(instanceType string) {
	// Load the shared AWS configuration
	cfg := utils.LoadAWSConfig()

	if instanceType == "" {
		fmt.Println("Instance type must be specified with -instance-type flag.")
		return
	}

	// Create an EC2 service client
	svc := utils.CreateEC2Client(cfg)

	// Describe instance types
	instanceTypes := []string{instanceType}
	instancesInfo := utils.DescribeInstanceTypes(svc, instanceTypes)

	// Create a table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Instance Type", "vCPUs", "Memory (MiB)"})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.FgHiGreenColor}, tablewriter.Colors{tablewriter.FgHiGreenColor}, tablewriter.Colors{tablewriter.FgHiGreenColor})
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	// Add instance type details to the table
	for _, info := range instancesInfo {
		table.Append([]string{info.InstanceType, fmt.Sprintf("%d", info.VCPUs), fmt.Sprintf("%d", info.Memory)})
	}

	// Render the table
	table.Render()
}
