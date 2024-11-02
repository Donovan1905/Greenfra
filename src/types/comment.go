package types

import (
	"github.com/fatih/color"
	"log"
	"strings"
)

type ResourceMetadata struct {
	ResourceReference string            `json:"resource_reference"`
	Metadata          map[string]string `json:"metadata"`
}

var AllowedResourceKeys = map[string][]string{
	"aws_instance":        {"usage_percentage"},
	"aws_lambda_function": {"monthly_invocation", "mean_execution_time"},
}

func ValidateComments(resources map[string]ResourceMetadata) error {
	var errors []string

	for key, resource := range resources {
		parts := strings.SplitN(key, ".", 2)
		if len(parts) != 2 {
			log.Fatalf("invalid resource key format in %s: must be 'resource_type.resource_name'", key)
		}

		resourceType := parts[0]

		allowedKeys, exists := AllowedResourceKeys[resourceType]
		if !exists {
			color.New(color.FgHiRed).Printf("comment for resource type %s are not handled for now\n", resourceType)
			errors = append(errors, key)
			continue
		}

		for metadataKey := range resource.Metadata {
			isAllowed := false
			for _, allowedKey := range allowedKeys {
				if allowedKey == metadataKey {
					isAllowed = true
					break
				}
			}
			if !isAllowed {
				color.New(color.FgHiRed).Printf("invalid metadata key '%s' for resource type '%s'\n", metadataKey, resourceType)
				errors = append(errors, key)
			}
		}
	}

	if len(errors) > 0 {
		log.Fatalf("Key Error when parsing greenfra comments")
	}

	return nil
}
