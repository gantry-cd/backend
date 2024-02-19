package models

type HomeResponse struct {
	OrganizationInfos []OrganizationInfos `json:"organizationInfos"`
}

type OrganizationInfos struct {
	Organization string   `json:"organization"`
	Repositories []string `json:"repositories"`
}
