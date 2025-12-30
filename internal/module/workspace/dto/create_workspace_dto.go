package workspace

type CreateWorkspaceDto struct {
	Name string  `json:"name" form:"name" validate:"required,min=3,max=195"`
	Slug *string `json:"slug,omitempty" form:"slug" validate:"omitempty"`
}
