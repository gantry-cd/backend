package models

type HomeResponse struct {
	OrganizationInfos []OrganizationInfos `json:"organizationInfos"`
}

type OrganizationInfos struct {
	Organization string   `json:"organization"`
	Repositories []string `json:"repositories"`
}

type GetRepositoryAppsRequest struct {
	Organization string `json:"organization"`
}

type GetRepositoryAppsResponse struct {
	Repositories []Repositories `json:"repositories"`
	Apps         []Apps         `json:"apps"`
}

type GetRepoBranchesRequest struct {
	Organization string `json:"organization"`
	Repository   string `json:"repository"`
}

type GetRepoBranchesResponse struct {
	Branches []Branches `json:"branches"`
}

type Branches struct {
	DeploymentName string `json:"deploymentName"`
	Branch         string `json:"branch"`
	PullRequestID  string `json:"pullRequestID"`
	Status         string `json:"status"`
	Version        string `json:"version"`
	Age            string `json:"age"`
}

type Repositories struct {
	Repository string `json:"repository"`
	Deployment int32  `json:"deployment"`
}

type Apps struct {
	AppName        string `json:"appName"`
	DeploymentName string `json:"deploymentName"`
	Status         string `json:"status"`
	Version        string `json:"version"`
	Age            string `json:"age"`
}
