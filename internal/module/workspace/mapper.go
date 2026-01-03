package workspace

import (
	"time"

	"github.com/google/uuid"
	"github.com/pimp13/jira-clone-backend-go/ent"
	entMembership "github.com/pimp13/jira-clone-backend-go/ent/membership"
	"github.com/pimp13/jira-clone-backend-go/internal/module/jwt"
)

type MembershipResponse struct {
	ID          uuid.UUID            `json:"id,omitempty"`
	Role        entMembership.Role   `json:"role,omitempty"`
	Status      entMembership.Status `json:"status,omitempty"`
	JoinedAt    time.Time            `json:"joined_at,omitempty"`
	UserID      uuid.UUID            `json:"user_id,omitempty"`
	WorkspaceID uuid.UUID            `json:"workspace_id,omitempty"`
	User        *jwt.UserInfo        `json:"user,omitempty"`
}

func ToMembershipResponse(m *ent.Membership) *MembershipResponse {
	var user *jwt.UserInfo
	if m.Edges.User != nil {
		user = jwt.ToUserInfo(m.Edges.User)
	}
	return &MembershipResponse{
		ID:          m.ID,
		Role:        m.Role,
		Status:      m.Status,
		JoinedAt:    m.JoinedAt,
		UserID:      m.UserID,
		WorkspaceID: m.WorkspaceID,
		User:        user,
	}
}

type WorkspaceResponse struct {
	ID          uuid.UUID             `json:"id,omitempty"`
	Name        string                `json:"name,omitempty"`
	Slug        string                `json:"slug,omitempty"`
	ImageURL    *string               `json:"imageUrl,omitempty"`
	ImagePath   *string               `json:"imagePath,omitempty"`
	InviteCode  string                `json:"inviteCode,omitempty"`
	CreatedAt   time.Time             `json:"createdAt,omitempty"`
	UpdatedAt   time.Time             `json:"updatedAt,omitempty"`
	Memberships []*MembershipResponse `json:"memberships,omitempty"`
}

func ToWorkspaceResponse(ws *ent.Workspace) *WorkspaceResponse {
	var memberships []*MembershipResponse

	if ws.Edges.Memberships != nil {
		memberships = make([]*MembershipResponse, 0, len(ws.Edges.Memberships))

		for _, m := range ws.Edges.Memberships {
			memberships = append(memberships, ToMembershipResponse(m))
		}
	}

	return &WorkspaceResponse{
		ID:          ws.ID,
		Name:        ws.Name,
		Slug:        ws.Slug,
		ImageURL:    ws.ImageURL,
		ImagePath:   ws.ImagePath,
		InviteCode:  ws.InviteCode,
		Memberships: memberships,
		CreatedAt:   ws.CreatedAt,
		UpdatedAt:   ws.UpdatedAt,
	}
}
