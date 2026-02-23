package model

type Touchpoint struct {
	ID             string   `json:"id"`
	Date           string   `json:"date"`
	Description    string   `json:"description"`
	Category       string   `json:"category"`
	Tags           []string `json:"tags"`
	PeopleInvolved []string `json:"people_involved"`
	URL            string   `json:"url"`
}

type TouchpointInput struct {
	Description    string   `json:"description"`
	Category       string   `json:"category"`
	Tags           []string `json:"tags"`
	PeopleInvolved []string `json:"people_involved"`
	URL            string   `json:"url"`
}

type Metadata struct {
	Categories []string `json:"categories"`
	Tags       []string `json:"tags"`
}
