package workspace

import (
	"github.com/pimp13/jira-clone-backend-go/ent"
	"github.com/pimp13/jira-clone-backend-go/internal/module/jwt"
)

func ToWorkspaceResponse(ws *ent.Workspace) *WorkspaceResponse {
	var owner *jwt.UserInfo

	if ws.Edges.Owner != nil {
		owner = jwt.ToUserInfo(ws.Edges.Owner)
	}

	return &WorkspaceResponse{
		ID:        ws.ID,
		Name:      ws.Name,
		Slug:      ws.Slug,
		ImageURL:  ws.ImageURL,
		OwnerID:   ws.OwnerID,
		CreatedAt: ws.CreatedAt,
		UpdatedAt: ws.UpdatedAt,
		Owner:     owner,
	}
}
