package models

type Package struct {
	Type          string                 `json:"type"`
	ParameterData map[string]interface{} `json:"parameter_data"`
	Outputs       map[string]interface{} `json:"outputs"`
}

type DeploymentMessage struct {
	ProjectID          string                 `json:"project_id"`
	PackageID          string                 `json:"package_id"`
	Package            Package                `json:"package"`
	ConnectedInputData map[string]interface{} `json:"connected_input_data"`
}
