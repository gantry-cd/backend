package models

type Pod struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Age    string `json:"age"`
	Image  string `json:"image"`
}

type BranchInfomationResponse struct {
	BranchName string `json:"branchName"`
	GitHubLink string `json:"gitHubLink"`
	Pods       []Pod  `json:"pods"`
	Yaml       string `json:"yaml"`
}
