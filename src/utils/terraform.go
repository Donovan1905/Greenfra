package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	fmt.Println("Terraform plan executed successfully.\n")
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
