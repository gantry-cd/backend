package models

type UsageRequest struct {
	Organization string `json:"organization"`
	Repository   string `json:"repository"`
	Span         int    `json:"span"`
}

type UsageResponse struct {
	Usages    []Usage `json:"usages"`
	IsDisable bool    `json:"isDisable"`
}

type Usage struct {
	PodName string `json:"podName"`
	CPU     string `json:"cpu"`
	MEM     string `json:"memory"`
	Storage string `json:"storage"`
}
