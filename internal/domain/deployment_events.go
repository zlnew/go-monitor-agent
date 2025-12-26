package domain

type EventDeploymentCreated struct {
	DeploymentID  int64 `json:"deployment_id"`
	ApplicationID int64 `json:"application_id"`
}

type EventDeploymentStatusChanged struct {
	DeploymentID  int64            `json:"deployment_id"`
	ApplicationID int64            `json:"application_id"`
	Status        DeploymentStatus `json:"status"`
}

type EventDeploymentCommitInfoReceived struct {
	DeploymentID  int64  `json:"deployment_id"`
	ApplicationID int64  `json:"application_id"`
	CommitHash    string `json:"commit_hash"`
	CommitMessage string `json:"commit_message"`
}
