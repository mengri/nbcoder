package clonepool

type Repository struct {
	ID    string `json:"id"`
	URL   string `json:"url"`
	Name  string `json:"name"`
	Branch string `json:"branch,omitempty"`
}
