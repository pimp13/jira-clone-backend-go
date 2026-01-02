package workspace

type UpdateWorkspaceDto struct {
	Name *string `json:"name,omitempty" form:"name" validate:"omitempty,min=3,max=195"`
	Slug *string `json:"slug,omitempty" form:"slug" validate:"omitempty"`
}
