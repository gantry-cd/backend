package models

type HomeResponse struct {
	OrganizationInfos []OrganizationInfos `json:"organizationInfos"`
}

type GetRepositoryAppsRequest struct {
	Organization string `json:"organization"`
}

type GetRepositoryAppsResponse struct {
	Repositories []Repositories `json:"repositories"`
	Apps         []Apps         `json:"apps"`
}

type GetBranchInfoRequest struct {
	Organization string `json:"organization"`
	Repository   string `json:"repository"`
}

type GetBranchInfoResponse struct {
	YAML string `json:"yaml"`
}

type OrganizationInfos struct {
	Organization string   `json:"organization"`
	Repositories []string `json:"repositories"`
}

type Repositories struct {
	Repository string `json:"repository"`
	Deployment int32  `json:"deployment"`
}

type Apps struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Version string `json:"version"`
	Age     string `json:"age"`
}
