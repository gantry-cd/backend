package models

type Repo struct {
	Name          string `json:"name"`
	PullRequestID string `json:"pullRequestID"`
	Branch        string `json:"branch"`
}

type App struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Version string `json:"version"`
	Image   string `json:"image"`
	Age     string `json:"age"`
}

type Organization struct {
	Name  string  `json:"name"`
	Repos []*Repo `json:"repos"`
	Apps  []*App  `json:"apps"`
}
