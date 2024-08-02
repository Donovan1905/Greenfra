package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"greenfra/src/services"
	"os/exec"
)

func ExecuteTerraformPlan() error {
	cmdPlan := exec.Command("terraform", "plan", "-out=tfplan")
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

func ExecuteTerraformShow() (map[string]interface{}, error) {
	cmdShow := exec.Command("terraform", "show", "-json", "tfplan")
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
