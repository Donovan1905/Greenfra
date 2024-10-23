package types

type ResourceMetadata struct {
	ResourceReference string            `json:"resource_reference"`
	Metadata          map[string]string `json:"metadata"`
}
