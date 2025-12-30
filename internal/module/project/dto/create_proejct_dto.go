package project

import "github.com/google/uuid"

type CreateProjectDto struct {
	Name        string    `json:"name" form:"name"`
	WorkspaceId uuid.UUID `json:"workspaceId" form:"workspaceId"`
	Description *string   `json:"description,omitempty" form:"description"`
	IsActive    string    `json:"isActive" form:"isActive"`
}
