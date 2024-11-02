package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"greenfra/src/services"
	"greenfra/src/types"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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

func ParseMetadataComments(tfFilePath string) (map[string]types.ResourceMetadata, error) {
	content, err := os.ReadFile(tfFilePath)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`/\*\s*greenfra\s*\n((?:[^\n]*\n)*?)\*/\s*resource\s+"(\w+)"\s+"(\w+)"\s*{`)
	matches := re.FindAllStringSubmatch(string(content), -1)

	resources := make(map[string]types.ResourceMetadata)

	for _, match := range matches {
		metadata := make(map[string]string)
		keyValuePairs := match[1]

		for _, line := range strings.Split(keyValuePairs, "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				metadata[key] = value
			}
		}

		resourceReference := fmt.Sprintf("%s.%s", match[2], match[3])
		resources[resourceReference] = types.ResourceMetadata{
			ResourceReference: resourceReference,
			Metadata:          metadata,
		}
	}

	return resources, nil
}

func ParseTfFilesInDirectory(dir string) (map[string]types.ResourceMetadata, error) {
	allResources := make(map[string]types.ResourceMetadata)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".tf") {
			resources, err := ParseMetadataComments(path)
			if err != nil {
				return fmt.Errorf("error parsing %s: %v", path, err)
			}
			for k, v := range resources {
				allResources[k] = v
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	err = types.ValidateComments(allResources)
	if err != nil {
		return nil, err
	}

	return allResources, nil
}
