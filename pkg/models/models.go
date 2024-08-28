package models

type Package struct {
	Type          string                 `json:"type"`
	ParameterData map[string]interface{} `json:"parameter_data"`
	Outputs       map[string]interface{} `json:"outputs"`
}

type DeploymentAction string

const (
	ActionDeploy  DeploymentAction = "DEPLOY"
	ActionDestroy DeploymentAction = "DESTROY"
)

type DeployStatus string

const (
	StartDeploy  DeployStatus = "DEPLOYING"
	StartDestroy DeployStatus = "DESTROYING"
	Deployed     DeployStatus = "DEPLOYED"
	Destroyed    DeployStatus = "NOT_DEPLOYED"
	Failed       DeployStatus = "FAILED"
)

type DeploymentMessage struct {
	ProjectID string  `json:"project_id"`
	PackageID string  `json:"package_id"`
	Package   Package `json:"package"`
	// i think this is actually the parameter data
	// and we still need to add the connections
	ConnectedInputData map[string]interface{} `json:"connected_input_data"`
	Action             DeploymentAction       `json:"action"`
	Secrets            map[string]string      `json:"secrets"`
}
