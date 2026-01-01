package workspace

import (
	"time"

	"github.com/google/uuid"
	"github.com/pimp13/jira-clone-backend-go/internal/module/jwt"
)

type WorkspaceResponse struct {
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Slug holds the value of the "slug" field.
	Slug string `json:"slug,omitempty"`
	// ImageURL holds the value of the "image_url" field.
	ImageURL *string `json:"imageUrl,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	// OwnerID holds the value of the "owner_id" field.
	OwnerID uuid.UUID `json:"ownerId,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	Owner *jwt.UserInfo `json:"owner,omitempty"`
}

type CreateWorkspaceResponse struct {
	ID uuid.UUID `json:"id,omitempty"`
}
