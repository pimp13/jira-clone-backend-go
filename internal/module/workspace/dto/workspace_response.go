package workspace

import (
	"github.com/google/uuid"
)

type CreateWorkspaceResponse struct {
	ID uuid.UUID `json:"id,omitempty"`
}

type UpdateWorkspaceResponse struct {
	ID uuid.UUID `json:"id,omitempty"`
}
