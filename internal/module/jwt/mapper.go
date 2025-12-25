package jwt

import "github.com/pimp13/jira-clone-backend-go/ent"

func ToUserInfo(u *ent.User) *UserInfo {
	return &UserInfo{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		IsActive:  u.IsActive,
		AvatarURL: u.AvatarURL,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
