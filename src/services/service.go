package services

type ResourceChange struct {
	Address string                 `json:"address"`
	Mode    string                 `json:"mode"`
	Type    string                 `json:"type"`
	Name    string                 `json:"name"`
	Change  map[string]interface{} `json:"change"`
}

type Service interface {
	Analyze(changes []ResourceChange) error
}
