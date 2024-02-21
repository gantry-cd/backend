package models

type PodLogRequest struct {
	Organization string `json:"organization"`
	Repository   string `json:"repository"`
	Pull         string `json:"pull"`
	Pod          string `json:"pod"`
}
