package models

type UsageRequest struct {
	Organization   string `json:"organization"`
	DeploymentName string `json:"deploymentName"`
}

type UsageResponse struct {
	Organization   string  `json:"organization"`
	DeploymentName string  `json:"deploymentName"`
	Usages         []Usage `json:"usages"`
	IsDisable      bool    `json:"isDisable"`
}

type Usage struct {
	PodName string `json:"podName"`
	CPU     int64  `json:"cpu"`
	MEM     int64  `json:"memory"`
	Storage int64  `json:"storage"`
}
