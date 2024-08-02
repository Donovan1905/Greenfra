package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"greenfra/src/services"
	"os/exec"
)

func ExecuteTerraformPlan(planPath string) error {
	cmdPlan := exec.Command("terraform", "plan", "-out", planPath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmdPlan.Stdout = &out
	cmdPlan.Stderr = &stderr
	err := cmdPlan.Run()
	if err != nil {
		return fmt.Errorf("terraform plan failed: %v\n%s", err, stderr.String())
	}
	return nil
}

func ExecuteTerraformShow(planPath string) (map[string]interface{}, error) {
	cmdShow := exec.Command("terraform", "show", "-json", planPath)
	var outShow bytes.Buffer
	var stderrShow bytes.Buffer
	cmdShow.Stdout = &outShow
	cmdShow.Stderr = &stderrShow
	err := cmdShow.Run()
	if err != nil {
		return nil, fmt.Errorf("terraform show failed: %v\n%s", err, stderrShow.String())
	}

	var result map[string]interface{}
	err = json.Unmarshal(outShow.Bytes(), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	cmdDeletePlan := exec.Command("rm", planPath)
	delErr := cmdDeletePlan.Run()
	if delErr != nil {
		return nil, fmt.Errorf("Failed to delete generated plan file")
	}

	return result, nil
}

func ExtractResourceChanges(plan map[string]interface{}) ([]services.ResourceChange, error) {
	var changes []services.ResourceChange

	if resourceChanges, ok := plan["resource_changes"].([]interface{}); ok {
		for _, change := range resourceChanges {
			var resourceChange services.ResourceChange
			changeBytes, err := json.Marshal(change)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(changeBytes, &resourceChange)
			if err != nil {
				return nil, err
			}
			changes = append(changes, resourceChange)
		}
	}

	return changes, nil
}
