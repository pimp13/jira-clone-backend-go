package workspace

import (
	"time"

	"github.com/google/uuid"
	"github.com/pimp13/jira-clone-backend-go/internal/module/jwt"
)

type WorkspaceResponse struct {
	ID        uuid.UUID     `json:"id,omitempty"`
	Name      string        `json:"name,omitempty"`
	Slug      string        `json:"slug,omitempty"`
	ImageURL  *string       `json:"imageUrl,omitempty"`
	CreatedAt time.Time     `json:"createdAt,omitempty"`
	UpdatedAt time.Time     `json:"updatedAt,omitempty"`
	OwnerID   uuid.UUID     `json:"ownerId,omitempty"`
	Owner     *jwt.UserInfo `json:"owner,omitempty"`
}

type CreateWorkspaceResponse struct {
	ID uuid.UUID `json:"id,omitempty"`
}

type UpdateWorkspaceResponse struct {
	ID uuid.UUID `json:"id,omitempty"`
}
