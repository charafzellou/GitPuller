package main

type Branch struct {
	Name string `json:"name"`
}

type Commit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Author struct {
			Date string `json:"date"`
		} `json:"author"`
	} `json:"commit"`
}

type Repository struct {
	Name    string `json:"name"`
	URL     string `json:"clone_url"`
	Private bool   `json:"private"`
}
