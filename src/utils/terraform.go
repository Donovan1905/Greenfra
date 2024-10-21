package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"greenfra/src/services"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func GetAWSRegion(tfplan map[string]interface{}) (string, error) {
	region := tfplan["configuration"].(map[string]interface{})["provider_config"].(map[string]interface{})["aws"].(map[string]interface{})["expressions"].(map[string]interface{})["region"].(map[string]interface{})["constant_value"].(string)

	return region, nil
}

func ParseMetadataComments(tfFilePath string) (map[string]string, string, string, error) {
	content, err := os.ReadFile(tfFilePath)
	if err != nil {
		return nil, "", "", err
	}

	lines := strings.Split(string(content), "\n")
	metadata := make(map[string]string)
	var resourceType, resourceIdentifier string

	for i, line := range lines {
		if strings.HasPrefix(line, "/* greenfra") {
			for j := i + 1; j < len(lines); j++ {
				line = lines[j]
				if strings.HasPrefix(line, "*/") {
					break
				}
				if strings.Contains(line, "=") {
					parts := strings.SplitN(line, "=", 2)
					if len(parts) == 2 {
						key := strings.TrimSpace(parts[0])
						value := strings.TrimSpace(parts[1])
						metadata[key] = value
					}
				}
			}
		}

		if strings.HasPrefix(line, "resource ") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				resourceType = parts[1]
				resourceIdentifier = parts[2]
			}
		}
	}

	return metadata, resourceType, resourceIdentifier, nil
}

func ParseTfFilesInDirectory(dir string) (map[string]map[string]string, error) {
	metadataMap := make(map[string]map[string]string)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".tf") {
			metadata, resourceType, resourceIdentifier, err := ParseMetadataComments(path)
			if err != nil {
				return fmt.Errorf("error parsing %s: %v", path, err)
			}

			if resourceIdentifier != "" && resourceType != "" && len(metadata) > 0 {
				metadataMap[resourceIdentifier] = metadata
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return metadataMap, nil
}
