package workspace

import (
	"time"

	"github.com/google/uuid"
	entMembership "github.com/pimp13/jira-clone-backend-go/ent/membership"
)

type MembershipResponse struct {
	ID          uuid.UUID            `json:"id,omitempty"`
	Role        entMembership.Role   `json:"role,omitempty"`
	Status      entMembership.Status `json:"status,omitempty"`
	JoinedAt    time.Time            `json:"joined_at,omitempty"`
	UserID      uuid.UUID            `json:"user_id,omitempty"`
	WorkspaceID uuid.UUID            `json:"workspace_id,omitempty"`
}

// func ToWorkspaceResponse(ws *ent.Workspace) *workspace.WorkspaceResponse {
// 	var owner *MembershipResponse

// 	if ws.Edges.Memberships != nil {
// 		owner = jwt.ToUserInfo(ws.Edges.Owner)
// 	}

// 	return &workspace.WorkspaceResponse{
// 		ID:        ws.ID,
// 		Name:      ws.Name,
// 		Slug:      ws.Slug,
// 		ImageURL:  ws.ImageURL,
// 		OwnerID:   ws.OwnerID,
// 		CreatedAt: ws.CreatedAt,
// 		UpdatedAt: ws.UpdatedAt,
// 		Owner:     owner,
// 	}
// }
