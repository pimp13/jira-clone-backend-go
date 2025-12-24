package workspace

type CreateWorkspaceDto struct {
	Name string  `json:"name"`
	Slug *string `json:"slug,omitempty"`
}
