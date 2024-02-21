package models

type UsageRequest struct {
	Organization string `json:"organization"`
	Repository   string `json:"repository"`
}

type UsageResponse struct {
	Usages    []Usage `json:"usages"`
	IsDisable bool    `json:"isDisable"`
}

type Usage struct {
	PodName string `json:"podName"`
	Branch  string `json:"branch"`
	PrID    string `json:"pullRequestID"`
	CPU     int64  `json:"cpu"`
	MEM     int64  `json:"memory"`
	Storage int64  `json:"storage"`
}
