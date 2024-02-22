package models

type PodLogRequest struct {
	Organization string `json:"organization"`
	Pod          string `json:"pod"`
}
